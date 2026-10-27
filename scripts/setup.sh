#!/usr/bin/env bash
# Installs and checks the tools make needs (16.8, B-07).
#
#   scripts/setup.sh [--profile translator|developer] [--check] [--yes]
#
# translator: git, make, Go (the e2e test server), Node, the npm deps of
#   frontend/ and e2e/, Playwright's Chromium and WebKit with their system
#   libraries. Enough for make i18n, make frontend, make screenshots.
# developer (default): all of that, plus the Wails CLI, golangci-lint in
#   .cache/bin, shellcheck, and the C toolchain with GTK and WebKitGTK
#   (Linux) or the Xcode Command Line Tools (macOS).
# --check installs nothing and only reports; make doctor runs it.
# --yes asks nothing.
#
# Versions: scripts/versions.env, go.mod (Go), e2e/package.json
# (Playwright). macOS uses Homebrew; Linux uses apt or dnf, and gets Go
# and Node from their official archives into ~/.local, without sudo.
# Every command is printed before it runs; sudo only after a yes.
# Bash 3.2 compatible: macOS ships that one.
set -euo pipefail

root=$(cd "$(dirname "$0")/.." && pwd)
# shellcheck source=scripts/versions.env
. "$root/scripts/versions.env"

usage() {
	sed -n '2,13p' "$0" | sed 's/^# \{0,1\}//'
	exit "${1:-0}"
}

profile=developer
check=0
yes=0
while [ $# -gt 0 ]; do
	case $1 in
	--profile) profile=${2:-}; shift ;;
	--profile=*) profile=${1#--profile=} ;;
	--check) check=1 ;;
	--yes | -y) yes=1 ;;
	-h | --help) usage 0 ;;
	*) echo "unknown option: $1" >&2; usage 2 >&2 ;;
	esac
	shift
done
case $profile in
translator | developer) ;;
*) echo "unknown profile: '$profile' (translator or developer)" >&2; exit 2 ;;
esac

go_want=$(awk '$1 == "go" { print $2; exit }' "$root/go.mod")
pw_want=$(sed -n 's/.*"@playwright\/test": *"\([^"]*\)".*/\1/p' "$root/e2e/package.json")
local_dir=$HOME/.local
golangci=${GOLANGCI:-}
if [ -z "$golangci" ] || [ "$golangci" = golangci-lint ]; then
	golangci='golangci-lint'
	[ -x "$root/.cache/bin/golangci-lint" ] && golangci=$root/.cache/bin/golangci-lint
fi

# --- Platform --------------------------------------------------------

os=$(uname -s)
machine=$(uname -m)
pm=
distro=$os
case $os in
Darwin)
	distro=macOS
	pm=brew
	;;
Linux)
	if [ -r /etc/os-release ]; then
		# shellcheck disable=SC1091
		distro=$(. /etc/os-release && echo "${PRETTY_NAME:-${ID:-Linux}}")
	fi
	if command -v apt-get > /dev/null 2>&1; then
		pm=apt
	elif command -v dnf > /dev/null 2>&1; then
		pm=dnf
	fi
	;;
esac
case $machine in
x86_64 | amd64) go_arch=amd64 node_arch=x64 ;;
arm64 | aarch64) go_arch=arm64 node_arch=arm64 ;;
*) go_arch='' node_arch='' ;;
esac

# --- Helpers ---------------------------------------------------------

# ver_ge a b: a >= b, for dotted numbers (1.27.2, 24).
ver_ge() {
	awk -v a="$1" -v b="$2" 'BEGIN {
		n = split(a, x, "."); m = split(b, y, ".")
		for (i = 1; i <= (n > m ? n : m); i++) {
			if (x[i] + 0 > y[i] + 0) exit 0
			if (x[i] + 0 < y[i] + 0) exit 1
		}
		exit 0
	}'
}

has() { command -v "$1" > /dev/null 2>&1; }

failed=0
todo=
hints=

# row tool found wanted status [required]: one line of the report.
# A required tool whose status is not OK fails the check.
row() {
	printf '  %-18s %-14s %-14s %s\n' "$1" "$2" "$3" "$4"
	if [ "$4" != OK ] && [ "${5:-required}" = required ]; then
		failed=1
	fi
}

need() { case " $todo " in *" $1 "*) ;; *) todo="$todo $1" ;; esac; }

hint() { hints="$hints
  $1"; }

# where tool dir...: sets found to the tool's path on PATH, else in one
# of the dirs, with onpath=0. Fails when it is nowhere.
where() {
	local tool=$1 dir
	shift
	onpath=1
	found=$(command -v "$tool" 2> /dev/null) && return 0
	onpath=0
	for dir in "$@"; do
		if [ -x "$dir/$tool" ]; then
			found=$dir/$tool
			return 0
		fi
	done
	return 1
}

# The step that brings git and make: the Xcode tools on macOS.
base_step() { if [ "$os" = Darwin ]; then echo xcode; else echo packages; fi; }

# --- Checks ----------------------------------------------------------

check_git() {
	if has git; then
		row git "$(git --version | awk '{ print $3 }')" any OK
	else
		row git - any missing
		need "$(base_step)"
	fi
}

check_make() {
	if has make; then
		row make "$(make --version 2> /dev/null | awk 'NR == 1 { print $NF }')" any OK
	else
		row make - any missing
		need "$(base_step)"
	fi
}

check_go() {
	local go v
	if ! where go "$local_dir/go/bin"; then
		row go - ">= $go_want" missing
		need go
		return
	fi
	go=$found
	# Inside the repo: go.mod's toolchain line may switch the version.
	v=$(cd "$root" && "$go" env GOVERSION 2> /dev/null | sed 's/^go//') || v=
	if [ -z "$v" ]; then
		row go "?" ">= $go_want" "wrong (go env fails)"
		need go
	elif ! ver_ge "$v" "$go_want"; then
		row go "$v" ">= $go_want" wrong
		need go
	elif [ "$onpath" = 0 ]; then
		row go "$v" ">= $go_want" "not on PATH"
		hint "add $(dirname "$go") to PATH"
	else
		row go "$v" ">= $go_want" OK
	fi
}

check_node() {
	local node v
	if ! where node "$local_dir/node/bin"; then
		row node - ">= $NODE_VERSION" missing
		need node
		return
	fi
	node=$found
	v=$("$node" --version | sed 's/^v//')
	if ! ver_ge "$v" "$NODE_VERSION"; then
		row node "$v" ">= $NODE_VERSION" wrong
		need node
	elif [ "$onpath" = 0 ]; then
		row node "$v" ">= $NODE_VERSION" "not on PATH"
		hint "add $(dirname "$node") to PATH"
	else
		row node "$v" ">= $NODE_VERSION" OK
	fi
}

# check_npm dir step: the dir's node_modules match its lock file.
check_npm() {
	if ! has npm; then
		row "npm deps $1" - package-lock missing
		need "$2"
	elif [ ! -d "$root/$1/node_modules" ]; then
		row "npm deps $1" - package-lock missing
		need "$2"
	elif ! (cd "$root/$1" && npm ls --depth=0 > /dev/null 2>&1); then
		row "npm deps $1" stale package-lock wrong
		need "$2"
	else
		row "npm deps $1" installed package-lock OK
	fi
}

check_playwright() {
	local f=$root/e2e/node_modules/@playwright/test/package.json v
	if [ ! -f "$f" ]; then
		row playwright - "$pw_want" missing
		need e2e-deps
		return
	fi
	v=$(sed -n 's/^ *"version": *"\([^"]*\)".*/\1/p' "$f" | head -n 1)
	if [ "$v" = "$pw_want" ]; then
		row playwright "$v" "$pw_want" OK
	else
		row playwright "$v" "$pw_want" wrong
		need e2e-deps
	fi
}

# check_browser name: every download Playwright lists for it is complete.
check_browser() {
	local out dir missing=0 rev
	if [ ! -f "$root/e2e/node_modules/@playwright/test/package.json" ] || ! has node; then
		row "$1" - "for $pw_want" missing
		need browsers
		return
	fi
	out=$(cd "$root/e2e" && npx --no-install playwright install --dry-run "$1" 2> /dev/null) || out=
	rev=$(echo "$out" | sed -n "s/.*(playwright $1 v\([0-9]*\)).*/\1/p" | head -n 1)
	for dir in $(echo "$out" | sed -n 's/^ *Install location: *//p'); do
		[ -f "$dir/INSTALLATION_COMPLETE" ] || missing=1
	done
	if [ -z "$out" ] || [ "$missing" = 1 ]; then
		row "$1" - "rev ${rev:-?}" missing
		need browsers
	else
		row "$1" "rev $rev" "rev $rev" OK
	fi
}

# check_browser_libs: the system libraries the browsers link (Linux).
check_browser_libs() {
	local ok=0
	case $pm in
	apt)
		[ -f "$root/e2e/node_modules/@playwright/test/package.json" ] && has node &&
			(cd "$root/e2e" && npx --no-install playwright install-deps --dry-run chromium webkit > /dev/null 2>&1) &&
			ok=1
		;;
	dnf)
		# shellcheck disable=SC2046
		rpm -q $(dnf_browser_pkgs) > /dev/null 2>&1 && ok=1
		;;
	*) return ;;
	esac
	if [ "$ok" = 1 ]; then
		row "browser libs" installed "$pm" OK
	else
		row "browser libs" - "$pm" missing
		need browser-libs
	fi
}

check_wails() {
	local wails v gobin
	gobin=$(go env GOPATH 2> /dev/null || echo "$HOME/go")/bin
	if ! where wails "$gobin"; then
		row wails - "$WAILS_VERSION" missing
		need wails
		return
	fi
	wails=$found
	v=$("$wails" version 2> /dev/null | grep -o '^v[0-9][0-9.]*' | head -n 1) || v=
	if [ "$v" != "$WAILS_VERSION" ]; then
		row wails "${v:-?}" "$WAILS_VERSION" wrong
		need wails
	elif [ "$onpath" = 0 ]; then
		row wails "$v" "$WAILS_VERSION" "not on PATH"
		hint "add $gobin to PATH"
	else
		row wails "$v" "$WAILS_VERSION" OK
	fi
}

# golangci-lint must match the version and be built with a Go at least
# as new as the repo's (docs/wiki/Linters.md).
check_golangci() {
	local out v built gov
	if ! has "$golangci"; then
		row golangci-lint - "$GOLANGCI_LINT_VERSION" missing
		need golangci
		return
	fi
	out=$("$golangci" version 2> /dev/null) || out=
	v=v$(echo "$out" | sed -n 's/.*has version v\{0,1\}\([0-9][0-9.]*\).*/\1/p')
	built=$(echo "$out" | sed -n 's/.*built with go\([0-9][0-9.]*\).*/\1/p')
	gov=$(cd "$root" && go env GOVERSION 2> /dev/null | sed 's/^go//') || gov=
	if [ "$v" != "$GOLANGCI_LINT_VERSION" ]; then
		row golangci-lint "$v" "$GOLANGCI_LINT_VERSION" wrong
		need golangci
	elif [ -n "$gov" ] && [ -n "$built" ] && ! ver_ge "$built" "$gov"; then
		row golangci-lint "$v" "$GOLANGCI_LINT_VERSION" "wrong (built with go$built < go$gov)"
		need golangci
	else
		row golangci-lint "$v" "$GOLANGCI_LINT_VERSION" OK
	fi
}

check_shellcheck() {
	if has shellcheck; then
		row shellcheck "$(shellcheck --version | awk '/^version:/ { print $2 }')" any OK optional
	else
		row shellcheck - any "missing (optional: make lint skips scripts)" optional
		if [ "$os" = Darwin ]; then need shellcheck; else need packages; fi
	fi
}

# The C side of the Wails build: cgo, GTK and WebKitGTK on Linux.
check_cbuild() {
	local pkg
	case $os in
	Darwin)
		if xcode-select -p > /dev/null 2>&1; then
			row "Xcode CLT" installed any OK
		else
			row "Xcode CLT" - any missing
			need xcode
		fi
		;;
	Linux)
		if has gcc; then
			row gcc "$(gcc -dumpversion)" any OK
		else
			row gcc - any missing
			need packages
		fi
		if ! has pkg-config; then
			row pkg-config - any missing
			need packages
			return
		fi
		row pkg-config "$(pkg-config --version)" any OK
		for pkg in gtk+-3.0 webkit2gtk-4.1; do
			if pkg-config --exists "$pkg"; then
				row "$pkg" "$(pkg-config --modversion "$pkg")" any OK
			else
				row "$pkg" - any missing
				need packages
			fi
		done
		;;
	esac
}

check_all() {
	failed=0
	todo=
	hints=
	echo "Profile: $profile; $distro ($machine), package manager: ${pm:-none}"
	printf '  %-18s %-14s %-14s %s\n' Tool Found Wanted Status
	check_git
	check_make
	check_go
	check_node
	check_npm frontend frontend-deps
	check_npm e2e e2e-deps
	check_playwright
	check_browser chromium
	if [ "$pm" = dnf ]; then
		# Its build links Ubuntu's libraries (ICU 74 and more).
		row webkit - - "skipped: Playwright's WebKit needs Debian or Ubuntu" optional
	else
		check_browser webkit
	fi
	if [ "$os" = Linux ]; then
		check_browser_libs
	fi
	if [ "$profile" = developer ]; then
		check_cbuild
		check_wails
		check_golangci
		check_shellcheck
	fi
	if [ "$os" = Linux ] && [ "$profile" = developer ]; then
		hint "Linux: make build, vet, test, lint and internal need TAGS=webkit2_41"
	fi
	if [ -n "$hints" ]; then
		echo "Notes:$hints"
	fi
}

# --- Installs --------------------------------------------------------

run() {
	printf '+ %s\n' "$*"
	"$@"
}

# as_root cmd...: root runs it, everyone else through sudo.
as_root() {
	if [ "$(id -u)" = 0 ]; then
		run "$@"
	elif has sudo; then
		run sudo "$@"
	else
		echo "needs root, and sudo is missing: $*" >&2
		exit 1
	fi
}

sudo_word() { [ "$(id -u)" = 0 ] || echo "sudo "; }

apt_pkgs() {
	echo "git make curl ca-certificates tar gzip"
	if [ "$profile" = developer ]; then
		echo "gcc pkg-config libgtk-3-dev libwebkit2gtk-4.1-dev shellcheck"
	fi
}

dnf_pkgs() {
	echo "git make curl ca-certificates tar gzip"
	if [ "$profile" = developer ]; then
		echo "gcc pkgconf-pkg-config gtk3-devel webkit2gtk4.1-devel ShellCheck"
	fi
}

# Chromium's libraries on Fedora; Playwright's install-deps knows only apt.
dnf_browser_pkgs() {
	echo "nss nspr atk at-spi2-atk at-spi2-core cups-libs libdrm libxkbcommon"
	echo "libXcomposite libXdamage libXfixes libXrandr libxshmfence mesa-libgbm"
	echo "pango cairo alsa-lib dbus-libs expat"
}

sha256() {
	if has sha256sum; then sha256sum "$1" | awk '{ print $1 }'; else shasum -a 256 "$1" | awk '{ print $1 }'; fi
}

# brew_get formula: install it, or upgrade the one brew already has.
brew_get() {
	if brew list --formula "$1" > /dev/null 2>&1; then
		run brew upgrade "$1"
	else
		run brew install "$1"
	fi
}

# fetch url file: download, failing on HTTP errors.
fetch() { run curl -fsSL --retry 3 -o "$2" "$1"; }

install_go_archive() {
	local file=go$go_want.linux-$go_arch.tar.gz tmp want
	tmp=$(mktemp -d)
	fetch "https://dl.google.com/go/$file" "$tmp/$file"
	want=$(curl -fsSL "https://dl.google.com/go/$file.sha256")
	if [ "$(sha256 "$tmp/$file")" != "$want" ]; then
		echo "checksum mismatch for $file" >&2
		exit 1
	fi
	run rm -rf "$local_dir/go"
	run mkdir -p "$local_dir"
	run tar -C "$local_dir" -xzf "$tmp/$file"
	rm -rf "$tmp"
}

install_node_archive() {
	local base=https://nodejs.org/dist/latest-v$NODE_VERSION.x tmp file sum
	tmp=$(mktemp -d)
	fetch "$base/SHASUMS256.txt" "$tmp/SHASUMS256.txt"
	file=$(awk -v s="-linux-$node_arch.tar.gz" 'index($2, s) && $2 ~ /^node-v/ { print $2; exit }' "$tmp/SHASUMS256.txt")
	sum=$(awk -v f="$file" '$2 == f { print $1 }' "$tmp/SHASUMS256.txt")
	fetch "$base/$file" "$tmp/$file"
	if [ -z "$file" ] || [ "$(sha256 "$tmp/$file")" != "$sum" ]; then
		echo "checksum mismatch for Node.js $file" >&2
		exit 1
	fi
	run rm -rf "$local_dir/node"
	run mkdir -p "$local_dir/node"
	run tar -C "$local_dir/node" --strip-components 1 -xzf "$tmp/$file"
	rm -rf "$tmp"
}

# describe step: one line of the plan.
describe() {
	case $1 in
	xcode) echo "Xcode Command Line Tools: run 'xcode-select --install' yourself, then this script again" ;;
	packages)
		case $pm in
		apt) echo "$(sudo_word)apt-get update; $(sudo_word)apt-get install -y $(apt_pkgs | tr '\n' ' ')" ;;
		dnf) echo "$(sudo_word)dnf install -y $(dnf_pkgs | tr '\n' ' ')" ;;
		esac
		;;
	go)
		case $pm in
		brew) echo "brew install go, or brew upgrade go (go.mod's toolchain line picks go$go_want)" ;;
		*) echo "Go $go_want from go.dev into $local_dir/go" ;;
		esac
		;;
	node)
		case $pm in
		brew) echo "brew install node, or brew upgrade node" ;;
		*) echo "Node.js $NODE_VERSION (latest) from nodejs.org into $local_dir/node" ;;
		esac
		;;
	shellcheck) echo "brew install shellcheck" ;;
	frontend-deps) echo "npm ci in frontend/" ;;
	e2e-deps) echo "npm ci in e2e/" ;;
	browsers)
		case $pm in
		dnf) echo "npx playwright install chromium (into Playwright's cache)" ;;
		*) echo "npx playwright install chromium webkit (into Playwright's cache)" ;;
		esac
		;;
	browser-libs)
		case $pm in
		apt) echo "the browsers' libraries: npx playwright install-deps chromium webkit (it runs sudo apt-get)" ;;
		dnf) echo "Chromium's libraries: $(sudo_word)dnf install -y $(dnf_browser_pkgs | tr '\n' ' ')" ;;
		esac
		;;
	wails) echo "go install github.com/wailsapp/wails/v2/cmd/wails@$WAILS_VERSION" ;;
	golangci) echo "go install golangci-lint@$GOLANGCI_LINT_VERSION into .cache/bin" ;;
	esac
}

install_step() {
	local gopath
	case $1 in
	xcode)
		echo "Install the Xcode Command Line Tools: xcode-select --install" >&2
		exit 1
		;;
	packages)
		case $pm in
		apt)
			as_root apt-get update
			# shellcheck disable=SC2046
			as_root env DEBIAN_FRONTEND=noninteractive apt-get install -y $(apt_pkgs)
			;;
		dnf)
			# shellcheck disable=SC2046
			as_root dnf install -y $(dnf_pkgs)
			;;
		esac
		;;
	shellcheck) brew_get shellcheck ;;
	go)
		case $pm in
		brew) brew_get go ;;
		*) install_go_archive ;;
		esac
		;;
	node)
		case $pm in
		brew) brew_get node ;;
		*) install_node_archive ;;
		esac
		;;
	frontend-deps) (cd "$root/frontend" && run npm ci --no-audit --no-fund) ;;
	e2e-deps) (cd "$root/e2e" && run npm ci --no-audit --no-fund) ;;
	browsers)
		case $pm in
		dnf) (cd "$root/e2e" && run npx --no-install playwright install chromium) ;;
		*) (cd "$root/e2e" && run npx --no-install playwright install chromium webkit) ;;
		esac
		;;
	browser-libs)
		case $pm in
		apt)
			if [ "$(id -u)" != 0 ]; then
				# It fails while anything is missing, after the list.
				echo "Playwright runs this with sudo:"
				(cd "$root/e2e" && npx --no-install playwright install-deps --dry-run chromium webkit) || true
			fi
			(cd "$root/e2e" && run npx --no-install playwright install-deps chromium webkit)
			;;
		dnf)
			# shellcheck disable=SC2046
			as_root dnf install -y $(dnf_browser_pkgs)
			;;
		esac
		;;
	wails) (cd "$root" && run go install "github.com/wailsapp/wails/v2/cmd/wails@$WAILS_VERSION") ;;
	golangci)
		(cd "$root" && GOBIN=$root/.cache/bin GOWORK=off run go install "github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$GOLANGCI_LINT_VERSION")
		golangci=$root/.cache/bin/golangci-lint
		;;
	esac
	# Later steps use what this one installed.
	gopath=$(go env GOPATH 2> /dev/null || echo "$HOME/go")
	PATH=$local_dir/go/bin:$local_dir/node/bin:$gopath/bin:$PATH
}

# --- Main ------------------------------------------------------------

old_path=$PATH
# Installs find what an earlier run put into ~/.local and GOPATH, even
# before the shell profile has it; --check reports it as not on PATH.
if [ "$check" = 0 ]; then
	PATH=$local_dir/go/bin:$local_dir/node/bin:$(go env GOPATH 2> /dev/null || echo "$HOME/go")/bin:$PATH
fi
if check_all; then :; fi
if [ "$check" = 1 ]; then
	if [ "$failed" = 1 ]; then
		echo "Missing or wrong tools: run scripts/setup.sh --profile $profile"
		exit 1
	fi
	echo "Everything the $profile profile needs is here."
	exit 0
fi

# The order matters: Node before npm, Go before go install.
plan=
for step in xcode packages shellcheck go node frontend-deps e2e-deps browsers browser-libs wails golangci; do
	case " $todo " in *" $step "*) plan="$plan $step" ;; esac
done
# Linux: curl and tar fetch Go and Node.
if [ "$os" = Linux ] && ! has curl; then
	case " $plan " in *" packages "*) ;; *) plan=" packages$plan" ;; esac
fi
if [ -z "$plan" ]; then
	echo "Nothing to install."
	exit "$failed"
fi
if [ -z "$pm" ] || [ -z "$go_arch" ]; then
	echo "Unsupported system: $distro ($machine)." >&2
	echo "This script knows macOS (Homebrew) and Linux with apt or dnf, on x86_64 or arm64." >&2
	echo "Install the tools above by hand: docs/wiki/Developer-setup.md." >&2
	exit 1
fi
if [ "$pm" = brew ] && ! has brew; then
	echo "Homebrew is missing. Install it from https://brew.sh, then run this script again." >&2
	exit 1
fi

echo
echo "The plan:"
for step in $plan; do
	echo "  - $(describe "$step")"
done
if [ "$yes" = 0 ]; then
	if [ ! -r /dev/tty ]; then
		echo "No terminal to ask; run with --yes to go ahead." >&2
		exit 1
	fi
	printf 'Go ahead? [y/N] '
	read -r answer < /dev/tty
	case $answer in
	y | Y | yes) ;;
	*) echo "Nothing installed."; exit 1 ;;
	esac
fi

for step in $plan; do
	echo
	echo "== $(describe "$step" | head -n 1)"
	install_step "$step"
done

echo
check_all
for dir in "$local_dir/go/bin" "$local_dir/node/bin" "$(go env GOPATH 2> /dev/null || echo "$HOME/go")/bin"; do
	case ":$old_path:" in
	*":$dir:"*) ;;
	*) [ -d "$dir" ] && echo "Add to your shell profile: export PATH=\"$dir:\$PATH\"" ;;
	esac
done
exit "$failed"
