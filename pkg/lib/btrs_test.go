package lib

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// The following samples were pulled from examples in the BTRS V3 guide
// docs/specifications/ANSI-X9-121201-BTRS-Format-Guide-Version-3.pdf
func TestBTRSParse(t *testing.T) {
	input := "01,122099999,123456789,150623,0200,1,,,3/"
	var header fileHeader
	_, err := header.parse(input)
	require.NoError(t, err)
	require.Equal(t, "122099999", header.Sender)
	require.Equal(t, "123456789", header.Receiver)
	require.Equal(t, "150623", header.FileCreatedDate)
	require.Equal(t, "0200", header.FileCreatedTime)
	require.Equal(t, "1", header.FileIdNumber)
	require.Equal(t, int64(0), header.PhysicalRecordLength)
	require.Equal(t, int64(0), header.BlockSize)
	require.Equal(t, int64(3), header.VersionNumber)
}
