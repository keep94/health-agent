package netif

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/Symantec/tricorder/go/tricorder/units"
	"os"
	"strings"
)

var filename string = "/proc/net/dev"

func (p *prober) probe() error {
	for _, netIf := range p.netInterfaces {
		netIf.probed = false
	}
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if err := p.processNetdevLine(scanner.Text()); err != nil {
			return err
		}
	}
	// TODO(rgooch): Clean up unprobed network interfaces once tricorder
	//               supports unregistration.
	return scanner.Err()
}

func (p *prober) processNetdevLine(line string) error {
	splitLine := strings.SplitN(line, ":", 2)
	if len(splitLine) != 2 {
		return nil
	}
	netIfName := strings.TrimSpace(splitLine[0])
	netIfData := splitLine[1]
	netIf := p.netInterfaces[netIfName]
	if netIf == nil {
		netIf = new(netInterface)
		p.netInterfaces[netIfName] = netIf
		metricsDir, err := p.dir.RegisterDirectory(netIfName)
		if err != nil {
			return err
		}
		if err := metricsDir.RegisterMetric("multicast-frames",
			&netIf.multicastFrames, units.None,
			"total multicast frames received or transmitted"); err != nil {
			return err
		}
		if err := metricsDir.RegisterMetric("rx-compressed-packets",
			&netIf.rxCompressedPackets, units.None,
			"compressed packets received"); err != nil {
			return err
		}
		if err := metricsDir.RegisterMetric("rx-data", &netIf.rxData,
			units.Byte, "bytes received"); err != nil {
			return err
		}
		if err := metricsDir.RegisterMetric("rx-dropped", &netIf.rxDropped,
			units.None, "receive packets dropped"); err != nil {
			return err
		}
		if err := metricsDir.RegisterMetric("rx-errors", &netIf.rxErrors,
			units.None, "total receive errors"); err != nil {
			return err
		}
		if err := metricsDir.RegisterMetric("rx-frame-errors",
			&netIf.rxFrameErrors, units.None,
			"receive framing errors"); err != nil {
			return err
		}
		if err := metricsDir.RegisterMetric("rx-overruns", &netIf.rxOverruns,
			units.None, "receive overrun errors"); err != nil {
			return err
		}
		if err := metricsDir.RegisterMetric("rx-packets", &netIf.rxPackets,
			units.None, "total packets received"); err != nil {
			return err
		}
		if err := metricsDir.RegisterMetric("tx-carrier-losses",
			&netIf.txCarrierLosses, units.None,
			"transmit carrier losses"); err != nil {
			return err
		}
		if err := metricsDir.RegisterMetric("tx-collision-errors",
			&netIf.txCollisionErrors, units.None,
			"transmit collision errors"); err != nil {
			return err
		}
		if err := metricsDir.RegisterMetric("tx-compressed-packets",
			&netIf.txCompressedPackets, units.None,
			"compressed packets transmitted"); err != nil {
			return err
		}
		if err := metricsDir.RegisterMetric("tx-data", &netIf.txData,
			units.Byte, "bytes transmitted"); err != nil {
			return err
		}
		if err := metricsDir.RegisterMetric("tx-dropped", &netIf.txDropped,
			units.None, "transmit packets dropped"); err != nil {
			return err
		}
		if err := metricsDir.RegisterMetric("tx-errors", &netIf.txErrors,
			units.None, "total transmit errors"); err != nil {
			return err
		}
		if err := metricsDir.RegisterMetric("tx-overruns", &netIf.txOverruns,
			units.None, "transmit overrun errors"); err != nil {
			return err
		}
		if err := metricsDir.RegisterMetric("tx-packets", &netIf.txPackets,
			units.None, "total packets transmitted"); err != nil {
			return err
		}
	}
	netIf.probed = true
	nScanned, err := fmt.Sscanf(netIfData,
		"%d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d",
		&netIf.rxData, &netIf.rxPackets, &netIf.rxErrors, &netIf.rxDropped,
		&netIf.rxOverruns, &netIf.rxFrameErrors, &netIf.rxCompressedPackets,
		&netIf.multicastFrames,
		&netIf.txData, &netIf.txPackets, &netIf.rxErrors, &netIf.txDropped,
		&netIf.txOverruns, &netIf.txCollisionErrors,
		&netIf.txCarrierLosses, &netIf.txCompressedPackets)
	if err != nil {
		return err
	}
	if nScanned < 16 {
		return errors.New(fmt.Sprintf("only read %d values from %s",
			nScanned, line))
	}
	return nil
}