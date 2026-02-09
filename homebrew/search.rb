# Homebrew formula for search CLI
# To install:
#   brew install jbutlerdev/tap/search

class Search < Formula
  desc "Command-line search tool using SearXNG"
  homepage "https://github.com/mule-ai/search"
  url "https://github.com/mule-ai/search/archive/refs/tags/v0.1.0.tar.gz"
  sha256 "SKIP"
  license "MIT"

  depends_on "go" => :build

  def install
    system "go", "build", *std_go_args(ldflags: "-s -w -X github.com/mule-ai/search/pkg/version.Version=#{version}")
  end

  test do
    version_output = shell_output("#{bin}/search --version")
    assert_match "search version #{version}", version_output
    assert_match "built with Go", version_output
  end
end