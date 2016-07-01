package cloudfoundry_test

import (
	"github.com/enaml-ops/enaml"
	. "github.com/enaml-ops/omg-cli/plugins/products/cloudfoundry/plugin"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("given a Diego Cell Partition", func() {
	Describe("given valid flags", func() {

		var instanceGroup *enaml.InstanceGroup
		var grouper InstanceGrouper

		Context("when ToInstanceGroup is called", func() {

			BeforeEach(func() {
				cf := new(Plugin)
				c := cf.GetContext([]string{
					"cloudfoundry",
					"--az", "eastprod-1",
					"--stemcell-name", "cool-ubuntu-animal",
					"--network", "foundry-net",
					"--diego-cell-ip", "10.0.0.39",
					"--diego-cell-ip", "10.0.0.40",
					"--diego-cell-vm-type", "cellvmtype",
					"--diego-cell-disk-type", "celldisktype",
				})
				grouper = NewDiegoCellPartition(c)
				instanceGroup = grouper.ToInstanceGroup()
			})

			It("then it should be populated with valid network configs", func() {
				ignet := instanceGroup.GetNetworkByName("foundry-net")
				Ω(ignet).ShouldNot(BeNil())
				Ω(ignet.StaticIPs).Should(ConsistOf("10.0.0.39", "10.0.0.40"))
			})

			It("then it should have an instance count in line with given IPs", func() {
				ignet := instanceGroup.GetNetworkByName("foundry-net")
				Ω(len(ignet.StaticIPs)).Should(Equal(instanceGroup.Instances))
			})

			It("then it should be populated the required jobs", func() {
				Ω(instanceGroup.GetJobByName("rep")).ShouldNot(BeNil())
				Ω(instanceGroup.GetJobByName("consul_agent")).ShouldNot(BeNil())
				Ω(instanceGroup.GetJobByName("cflinuxfs2-rootfs-setup")).ShouldNot(BeNil())
				Ω(instanceGroup.GetJobByName("garden")).ShouldNot(BeNil())
				Ω(instanceGroup.GetJobByName("statsd-injector")).ShouldNot(BeNil())
				Ω(instanceGroup.GetJobByName("metron_agent")).ShouldNot(BeNil())
			})

			Describe("given a rep job", func() {
				Context("when defined", func() {
					It("then it should use the correct release", func() {
						job := instanceGroup.GetJobByName("rep")
						Ω(job.Release).Should(Equal(DiegoReleaseName))
					})
				})
			})

			Describe("given a consul_agent job", func() {
				Context("when defined", func() {
					It("then it should use the correct release", func() {
						job := instanceGroup.GetJobByName("consul_agent")
						Ω(job.Release).Should(Equal(CFReleaseName))
					})
				})
			})

			Describe("given a cflinuxfs2-rootfs-setup job", func() {
				Context("when defined", func() {
					It("then it should use the correct release", func() {
						job := instanceGroup.GetJobByName("cflinuxfs2-rootfs-setup")
						Ω(job.Release).Should(Equal(CFLinuxFSReleaseName))
					})
				})
			})

			Describe("given a garden job", func() {
				Context("when defined", func() {
					It("then it should use the correct release", func() {
						job := instanceGroup.GetJobByName("garden")
						Ω(job.Release).Should(Equal(GardenReleaseName))
					})
				})
			})

			Describe("given a statsd-injector job", func() {
				Context("when defined", func() {
					It("then it should use the correct release", func() {
						job := instanceGroup.GetJobByName("statsd-injector")
						Ω(job.Release).Should(Equal(CFReleaseName))
					})
				})
			})

			Describe("given a metron_agent job", func() {
				Context("when defined", func() {
					It("then it should use the correct release", func() {
						job := instanceGroup.GetJobByName("metron_agent")
						Ω(job.Release).Should(Equal(CFReleaseName))
					})
				})
			})
		})
	})
})
