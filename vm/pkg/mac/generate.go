package mac

import (
	"crypto/rand"
	"fmt"
)

func Generate() (string, error) {
	var mac [6]byte

	// Randomly generate 6 bytes
	_, err := rand.Read(mac[:])
	if err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}

	// Set the locally administered address (LAA) bit (bit 1 of the first byte)
	// Ensure that the first byte has the multicast bit (bit 0) cleared and the locally administered bit (bit 1) set
	// This makes sure it's a valid locally administered address (LAA).
	mac[0] = (mac[0] & 0xFE) | 0x02

	// Return the MAC address as a string in the format "XX:XX:XX:XX:XX:XX"
	return fmt.Sprintf("%02X:%02X:%02X:%02X:%02X:%02X", mac[0], mac[1], mac[2], mac[3], mac[4], mac[5]), nil
}
