class Aryflow < Formula
  desc "CLI tool for AryFlow workflow automation"
  homepage "https://github.com/EslavaDev/aryflow"
  version "0.2.1"
  license "MIT"

  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/EslavaDev/aryflow/releases/download/v0.2.1/aryflow_0.2.1_darwin_arm64.tar.gz"
      sha256 "f498fdfef1d142246cba6772aaa9f3587df86a1f027e53d2f23b3c7479d94555"
    else
      url "https://github.com/EslavaDev/aryflow/releases/download/v0.2.1/aryflow_0.2.1_darwin_amd64.tar.gz"
      sha256 "54cca0a8d80c97e086f0652361b0ee2db63d109492ae721441f6400440e35d3f"
    end
  end

  on_linux do
    url "https://github.com/EslavaDev/aryflow/releases/download/v0.2.1/aryflow_0.2.1_linux_amd64.tar.gz"
    sha256 "9c27885f2800e5bc144809a93c1687491a8a7d86cec5456c72401642e5e244de"
  end

  def install
    bin.install "aryflow"
  end

  test do
    assert_match version.to_s, shell_output("#{bin}/aryflow --version")
  end
end
