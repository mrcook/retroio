package hardware

type Compatibility uint8

// Computer Hardware compatibility flags.
const (
	Runs         Compatibility = 0x00 // "The tape RUNS on this machine or with this hardware, but may or may not use the hardware or special features of the machine."
	RunsUsesHW   Compatibility = 0x01 // "The tape USES the hardware or special features of the machine, such as extra memory or a sound chip."
	RunsLegacyHW Compatibility = 0x02 // "The tape RUNS but it DOESN'T use the hardware or special features of the machine."
	NoRun        Compatibility = 0x03 // "The tape DOESN'T RUN on this machine or with this hardware."
)

// Information detailing the relationship between a piece of software and the hardware.
var CompatibilityDescriptions = map[Compatibility]string{
	Runs:         "The tape RUNS on this machine or with this hardware, but may or may not use the hardware or special features of the machine.",
	RunsUsesHW:   "The tape USES the hardware or special features of the machine, such as extra memory or a sound chip.",
	RunsLegacyHW: "The tape RUNS but it DOESN'T use the hardware or special features of the machine.",
	NoRun:        "The tape DOESN'T RUN on this machine or with this hardware.",
}
