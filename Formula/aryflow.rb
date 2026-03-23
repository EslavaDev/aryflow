class Aryflow < Formula
  desc "CLI tool for AryFlow workflow automation"
  homepage "https://github.com/EslavaDev/aryflow"
  version "0.2.0"
  license "MIT"

  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/EslavaDev/aryflow/releases/download/v0.2.0/aryflow_0.2.0_darwin_arm64.tar.gz"
      sha256 "b8994dfdf8c80f9cb126645ae23b3edc72319b2459a352960fb8e7c15fd091bd"
    else
      url "https://github.com/EslavaDev/aryflow/releases/download/v0.2.0/aryflow_0.2.0_darwin_amd64.tar.gz"
      sha256 "f1080f37cd68c796c0fb1b52e73306dccd2aa1ee57f644f31b8dfc66adf1062f"
    end
  end

  on_linux do
    url "https://github.com/EslavaDev/aryflow/releases/download/v0.2.0/aryflow_0.2.0_linux_amd64.tar.gz"
    sha256 "5cc1bbc7dd493e26bfa16a45e747863d11ce97be27c7f8a7850036f5ad3ebcd6"
  end

  def install
    bin.install "aryflow"
  end

  test do
    assert_match version.to_s, shell_output("#{bin}/aryflow --version")
  end
end
