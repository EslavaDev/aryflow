class Aryflow < Formula
  desc "CLI tool for AryFlow workflow automation"
  homepage "https://github.com/EslavaDev/aryflow"
  version "0.2.2"
  license "MIT"

  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/EslavaDev/aryflow/releases/download/v0.2.2/aryflow_0.2.2_darwin_arm64.tar.gz"
      sha256 "030168f961db0b7e693f6a2188a24d47a96085a5db64bc0d00f4353e74135c95"
    else
      url "https://github.com/EslavaDev/aryflow/releases/download/v0.2.2/aryflow_0.2.2_darwin_amd64.tar.gz"
      sha256 "2e97742e3006132b98883760bf6ca1abacc5448771b3150c8682857592311868"
    end
  end

  on_linux do
    url "https://github.com/EslavaDev/aryflow/releases/download/v0.2.2/aryflow_0.2.2_linux_amd64.tar.gz"
    sha256 "daa7248a6f35d59de39f45415fba5ed6e5ffba5cbeb308a15ed69426f4227844"
  end

  def install
    bin.install "aryflow"
  end

  test do
    assert_match version.to_s, shell_output("#{bin}/aryflow --version")
  end
end
