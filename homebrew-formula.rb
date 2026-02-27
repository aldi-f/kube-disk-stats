class Kubediskstats < Formula
  desc "CLI tool for querying Kubernetes node and pod disk usage statistics"
  homepage "https://github.com/aldi-f/kube-disk-stats"
  url "https://github.com/aldi-f/kube-disk-stats.git",
    tag:      "VERSION",
    revision: "REVISION"

  depends_on "go" => :build

  def install
    system "go", "build", *std_go_args(ldflags: "-s -w")
  end

  test do
    system bin/"kube-disk-stats", "version"
  end
end
