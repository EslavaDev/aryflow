class Aryflow < Formula
  desc "CLI tool for AryFlow workflow automation"
  homepage "https://github.com/EslavaDev/aryflow"
  version "0.2.3"
  license "MIT"

  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/EslavaDev/aryflow/releases/download/v0.2.3/aryflow_0.2.3_darwin_arm64.tar.gz"
      sha256 "77387252689d6714779ccfd1565ac926b5078e78ebb3794b320e611546ccc2db"
    else
      url "https://github.com/EslavaDev/aryflow/releases/download/v0.2.3/aryflow_0.2.3_darwin_amd64.tar.gz"
      sha256 "e8175e252e06ce0a9eff4494365d86907415026911191e5732a15dfa29aa93bd"
    end
  end

  on_linux do
    url "https://github.com/EslavaDev/aryflow/releases/download/v0.2.3/aryflow_0.2.3_linux_amd64.tar.gz"
    sha256 "f43733c24c63a048a67396c599c891e014d738d88a62966c884f2c249d557e14"
  end

  def install
    bin.install "aryflow"
  end

  test do
    assert_match version.to_s, shell_output("#{bin}/aryflow --version")
  end
end
