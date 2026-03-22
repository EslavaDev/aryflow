class Aryflow < Formula
  desc "CLI tool for AryFlow workflow automation"
  homepage "https://github.com/EslavaDev/aryflow"
  version "0.1.0"
  license "MIT"

  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/EslavaDev/aryflow/releases/download/v#{version}/aryflow_#{version}_darwin_arm64.tar.gz"
      sha256 "PLACEHOLDER"
    else
      url "https://github.com/EslavaDev/aryflow/releases/download/v#{version}/aryflow_#{version}_darwin_amd64.tar.gz"
      sha256 "PLACEHOLDER"
    end
  end

  on_linux do
    url "https://github.com/EslavaDev/aryflow/releases/download/v#{version}/aryflow_#{version}_linux_amd64.tar.gz"
    sha256 "PLACEHOLDER"
  end

  def install
    bin.install "aryflow"
  end

  test do
    assert_match version.to_s, shell_output("#{bin}/aryflow --version")
  end
end
