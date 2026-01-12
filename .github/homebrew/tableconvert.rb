class Tableconvert < Formula
  desc "Universal table format converter - CSV, JSON, MySQL, Markdown, Excel, and more"
  homepage "https://github.com/martianzhang/tableconvert"
  url "https://github.com/martianzhang/tableconvert/archive/refs/tags/v1.0.0-pre.tar.gz"
  # sha256 "TO_BE_CALCULATED" # Run: shasum -a 256 v1.0.0-pre.tar.gz
  license "MIT"
  head "https://github.com/martianzhang/tableconvert.git", branch: "main"

  depends_on "go" => :build

  def install
    # Build from source
    system "go", "build", *std_go_args(ldflags: "-s -w"), "./cmd/tableconvert"
  end

  test do
    # Test version flag
    assert_match "tableconvert version", shell_output("#{bin}/tableconvert --version")

    # Test basic conversion
    (testpath/"test.csv").write "name,age\nAlice,30\nBob,25"
    output = shell_output("#{bin}/tableconvert --from=csv --to=json test.csv")
    assert_match '"name": "Alice"', output
    assert_match '"age": "30"', output

    # Test help command
    assert_match "tableconvert", shell_output("#{bin}/tableconvert --help")

    # Test stdin/stdout conversion
    output = pipe_output("#{bin}/tableconvert --from=csv --to=json", "name,age\nAlice,30")
    assert_match '"name": "Alice"', output
  end
end
