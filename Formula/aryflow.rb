class Aryflow < Formula
  desc "CLI tool for AryFlow workflow automation"
  homepage "https://github.com/EslavaDev/aryflow"
  version "0.1.0"
  license "MIT"

  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/EslavaDev/aryflow/releases/download/v0.1.0/aryflow_0.1.0_darwin_arm64.tar.gz"
      sha256 "e13431c1ce3c9ae4a9e387ced69e60e12c807ad37c265624fb53b7714ecfa0fd"
    else
      url "https://github.com/EslavaDev/aryflow/releases/download/v0.1.0/aryflow_0.1.0_darwin_amd64.tar.gz"
      sha256 "af6ada0fb6c746e7347017f4fadb8eed63fcea366d262cb635f0392e9adb0d4b"
    end
  end

  on_linux do
    url "https://github.com/EslavaDev/aryflow/releases/download/v0.1.0/aryflow_0.1.0_linux_amd64.tar.gz"
    sha256 "51dd3e136d7f43eeb92f44519cf48e530dcdb67cb1166d88cc0f316c3593f5fa"
  end

  def install
    bin.install "aryflow"
  end

  test do
    assert_match version.to_s, shell_output("#{bin}/aryflow --version")
  end
end
