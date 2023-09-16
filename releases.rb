
DEFAULT_BRANCH="master"
RELEASE_NOTES="RELEASES.md"

WINDOWS=(RUBY_PLATFORM =~ /mswin|mingw|cygwin/)
$stdout.sync = true


# ----------------------------------------------------------------------------
# Define `memoized_attr` helper on Object to cache results of expensive
# operations in the subsequent class definitions
# ----------------------------------------------------------------------------

class Object
    def self.memoized_attr(attr_name, &block)
        @@memoized_attrs ||= []
        @@memoized_attrs << attr_name

        if block
            define_method("_#{attr_name}", &block)
            private "_#{attr_name}"
        end

        define_method(attr_name) do
            instance_variable_defined?("@#{attr_name}") ?
                instance_variable_get("@#{attr_name}") :
                instance_variable_set("@#{attr_name}", send("_#{attr_name}"))
        end
    end

    def reset_memoized_attrs
        @@memoized_attrs.each do |attr_name|
            remove_instance_variable("@#{attr_name}") if instance_variable_defined?("@#{attr_name}")
        end
    end
end


class Foo
    memoized_attr :name
    memoized_attr :name2 do
        puts "_name2()"
        name
    end


    private
    def _name()
        puts "_name()"
        "my-name"
    end
end

f = Foo.new



def git(cmd)
    puts "git #{cmd}"
    return `git #{cmd}`.strip()
end

GIT_LOG_SPLIT=/^(?:\((.*?)\))?\s*(.*?)\s*$/
LogEntry = Struct.new(:msg, :hash, :refs, :tags, :pr)

# git_log returns the commit log for either the local HEAD or the specified ref.
# Each commit is represented by a LogEntry containing the commit message, commit
# hash, and a list of refs and tags.
def git_log(ref='', limit: 100)
    logs = git( "log --pretty=oneline --decorate #{ref} --max-count=#{limit}" ).split("\n").map { |l|
        h, l = l.split(/\s+/, 2)
        m = GIT_LOG_SPLIT.match(l)
        rr, l = [m[1] || "", m[2]]
        tags, refs = rr.split(/,\s*/).partition { |r| r =~ /^tag:\s*(.*)/ }
        tags = tags.map { |t| t.sub(/^tag:\s*/, '')}
        LogEntry.new(l, h,refs, tags)
    }
end

# ----------------------------------------------------------------------------
# GitInfo : Helper to extract project history inforrmation from git repo
# ----------------------------------------------------------------------------

class GitInfo
    class << self
        def default() return @default ||= new end
    end

    def initialize()
        if git('rev-parse --is-shallow-repository') == 'true'
            puts "Fetching missing information from remote ..."
            system('git fetch --prune --tags --unshallow')
        end
    end


    memoized_attr :dir    do git('rev-parse --show-toplevel') end
    memoized_attr :commit do git('rev-parse HEAD') end
    memoized_attr :branch do git("rev-parse --abbrev-ref HEAD") end

    memoized_attr :name do
        remote_basename = File.basename(remote() || "" )
        return remote_basename if remote_basename != ""
        return File.basename(File.expand_path("."))
    end

    memoized_attr :remote do
        remote = git('remote get-url origin')
        m = GIT_SSH_REPO.match(remote)
        return remote if m.nil?

        host = m[:host]
        host = "github.com" if host.end_with? ("github.com")
        return "https://#{host}/#{m[:path]}/"
    end
    GIT_SSH_REPO = /git@(?<host>[^:]+):(?<path>.+)(?:.git)?/

    # def name()      return @name    ||= _name()     end
    # def version()   return @version ||= _version()  end
    # def mtag()      return @mtag    ||= _mtag()     end

end


# semver_parse parses a string containing a version string compatible with
# semantic versionning 2.0.0, with an optional 'v' prefix, and returns a array
# of 3 arrays containing the base version, the pre identifiers, and the build
# ientifiers. Parsing is lax and allows malformed version strings to be parsed
# anyway. The base version is always converted to integers or 0s, and truncated
# or 0-padded to exactly 3 integers. Pre and build fragments are returned as
# integers and strings; integers accept leading zeroes and strings are stripped
# of any invalid characters; empty fragments are dropped. When printed back out,
# the resulting string should be an exact match for valid semver or the
# "closest" valid semever for invalid strings.
def semver_parse(v)
    def fragment(v) v = v.gsub(/[^a-zA-Z0-9-]/, ''); Integer(v, exception: false) || (v if !v.empty?) end
    a,b = ((v[1..-1] if v[0] == 'v') || v).split('-', 2)
    b,c = (b || '').split('+', 2)
    a = (((a||'').split('.').map{|v| v.to_i}) + [0]*3)[0...3]
    b = (b||'').split('.').map{ |v| fragment(v) }.select { |v| !v.nil? }
    c = (c||'').split('.').map{ |v| fragment(v) }.select { |v| !v.nil? }
    return [a,b,c]
end

# semver_format formats the output of semver_parse back into a string, with a
# configurable prefix defaulting to 'v'.
def semver_format(v, prefix:'v')
    a, b, c = v
    a = (a.to_a + [0]*3)[0...3]
    v = prefix + a.join('.')
    v += '-' + b.join('.') if !b.nil? && !b.empty?
    v += '+' + c.join('.') if !c.nil? && !c.empty?
    v
end

# semver_comp compares 2 version strings according to semver 2.0.0 semantic and
# returns an output similar to <=>.
def semver_comp(a, b)
    def cmp(a,b) a <=> b || (a.nil? ? -1 : b.nil? ? 1 : a.is_a?(Numeric) ? -1 : 1) end
    def zip_max(a, b) a.length >= b.length ? a.zip(b) : b.zip(a).map(&:reverse) end

    a, b = semver_parse(a), semver_parse(b)
    c = a[0] <=> b[0]
    return c if c != 0
    return b[1].length <=> a[1].length if a[1].length == 0 || b[1].length == 0
    zip_max(a[1], b[1]).each { |aa,bb| c = cmp(aa,bb); return c if c != 0 }
    return 0
end


# ----------------------------------------------------------------------------
# BuildInfo : Helper to extract version inforrmation for git repo
# ----------------------------------------------------------------------------

class BuildInfo
    class << self
        def default() return @default ||= new end
    end

    def initialize()
        if git('rev-parse --is-shallow-repository') == 'true'
            puts "Fetching missing information from remote ..."
            system('git fetch --prune --tags --unshallow')
        end
    end

    def name()      return @name    ||= _name()     end
    def version()   return @version ||= _version()  end
    def remote()    return @remote  ||= _remote()   end
    def commit()    return @commit  ||= _commit()   end
    def dir()       return @dir     ||= _dir()      end
    def branch()    return @branch  ||= _branch()   end
    def mtag()      return @mtag    ||= _mtag()     end

    def on_release_branch(b, v)
        return b == DEFAULT_BRANCH || (!v.nil? && v.start_with?("#{b}."))
    end

    def reset()
        @name = @version = @remote = @commit = @dir = @branch = @mtag = nil
    end

    private
    def _commit()   return git('rev-parse HEAD')               end
    def _dir()      return git('rev-parse --show-toplevel')    end
    def _branch()   return git("rev-parse --abbrev-ref HEAD").strip.gsub(/[^A-Za-z0-9\._-]+/, '-') end

    def _name()
        remote_basename = File.basename(remote() || "" )
        return remote_basename if remote_basename != ""
        return File.basename(File.expand_path("."))
    end

    def _version()
        v, b, n, g = _info()                    # Extract base info from git branch and tags
        m = _mtag()                             # Detect locally modified files
        v = _patch(v) if n > 0 || !m.nil?       # Increment patch if needed to to preserve semver orderring
        b = 'rc' if on_release_branch(b, v)     # Rename branch fragment to 'rc' for default or release maintenance branch
        return v if b == 'rc' && n == 0 && m.nil?
        return "#{v}-" + [b, n, g, m].compact().join('.')
    end

    def _info()
        # Note: Due to glob(7) limitations, the following pattern enforces
        # 3-part dot-separated sequences starting with a digit,
        # rather than 3 dot-separated numbers.
        pattern = WINDOWS ? '"v[0-9]*.[0-9]*.[0-9]*"' : "'v[0-9]*.[0-9]*.[0-9]*'"
        d = git("describe --always --tags --long --match #{pattern}").strip.split('-')
        if d.count != 0
            return ['v0.0.0', branch, git("rev-list --count HEAD").strip.to_i, "g#{d[0]}"] if d.count == 1
            return [d[0], branch, d[1].to_i, d[2]] if d.count == 3
        end
        return ['v0.0.0', "none", 0, 'g0000000']
    end

    def _patch(v)
        # Increment the patch number by 1, so that intermediate version strings
        # sort between the last tag and the next tag according to semver.
        #   v0.6.1
        #       v0.6.1-maa-cleanup.1.g6ede8cd   <-- with _patch()
        #   v0.6.0
        #       v0.6.0-maa-cleanup.1.g6ede8cd   <-- without _patch()
        #   v0.5.99
        vv = v[1..-1].split('.').map { |v| v.to_i }
        vv[-1] += 1
        v = "v" + vv.join(".")
        return v
    end

    def _mtag()
        # Generate a `.mXXXXXXXX` fragment based on latest mtime of modified
        # files in the index. Returns `nil` if no files are locally modified.
        status = git("status --porcelain=2 --untracked-files=no")
        files = status.lines.map {|l| l.strip.split(/ +/).last }.map { |n| n.split(/\t/).first }
        t = files.map { |f| File.mtime(f).to_i rescue nil }.compact.max
        return t.nil? ? nil : "m%08x" % t
    end

    GIT_SSH_REPO = /git@(?<host>[^:]+):(?<path>.+).git/
    def _remote()
        remote = git('remote get-url origin')
        m = GIT_SSH_REPO.match(remote)
        return remote if m.nil?

        host = m[:host]
        host = "github.com" if host.end_with? ("github.com")
        return "https://#{host}/#{m[:path]}/"
    end
end

# ----------------------------------------------------------------------------

class ReleaseInfo
    class << self
        def default() return @default ||= new end
    end

    def initialize(release_notes_file = RELEASE_NOTES)
        @release_notes_file = release_notes_file
    end

    def reset()
        @loaded = false
        return self
    end

    # def release()           @loaded ||= _load(); return "_release()" end
    # def release_version()   @loaded ||= _load(); return "_release_version()" end
    # def release_branch()    @loaded ||= _load(); return "_release_branch()" end
    # def release_notes()     @loaded ||= _load(); return "_release_notes()" end
    # def no_release_reason() @loaded ||= _load(); return "_no_release_reason()" end

    def versions()          @loaded ||= _load(); return @versions end
    def commits()           @loaded ||= _load(); return @commits end

    private
    def _load()
        puts "Loading..."
        @release_notes = _load_release_note(@release_notes_file)
        @versions = _sort_versions(@release_notes.keys())
        @commits = git_log()

        return true
    end

    def _load_release_note(filename)
        version = nil
        return File.readlines(filename).chunk do |l|
            version = l[2..-1].strip() if l.start_with?( "# "); version
        end.map do |v, ll|
            ll.shift
            ll.shift while ll.first.strip == ""
            ll.pop while ll.last.strip == ""
            [v, ll.map { |l| l.rstrip }]
        end.reverse.to_h
    end

    def _sort_versions(versions)
        versions.sort {|a,b| semver_comp(b, a) }
    end
end

# ----------------------------------------------------------------------------

# ReleaseInfo.default.versions
# ReleaseInfo.default.commits

# def rsplit(s, sep, n=-1)
#     v = s.split(sep)
#     if n > 0
#         v = [v[0...-(n-1)].join(sep)] + v[-(n-1)..-1]
#     end
#     v
# end

# commits.each_with_index.flat_map{|c,i| c.tags.map {|t| [t,i]}}.select {|t,i| t=~/^v?\d+\.\d+\.\d+$/}.first
