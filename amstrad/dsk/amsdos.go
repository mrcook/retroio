// Amstrad CP/M Disc reading.
package dsk

import (
	"bytes"
	"encoding/binary"
	"io"

	"github.com/pkg/errors"

	"retroio/amstrad/dsk/amsdos"
)

type AmsDos struct {
	DPB         amsdos.DiskParameterBlock
	Directories []amsdos.Directory
}

// Read the contents of an AMSDOS formatted disk
func (a *AmsDos) Read(disk *DSK) error {
	if len(disk.Tracks) == 0 {
		return errors.New("no available tracks")
	}
	track := disk.Tracks[0]

	if len(track.Sectors) == 0 {
		return errors.New("no sectors found")
	}

	sectorSize, ok := sectorSizeMap[track.SectorSize]
	if !ok {
		return errors.Errorf("invalid sector size: 0x%02X", track.SectorSize)
	}

	a.readDirectories(sectorSize, &track)

	// must be executed after reading the directories
	a.generateDPB(disk.Info.TrackSize, sectorSize, track.Sectors[0].ID, disk.Info.mediaType())

	return nil
}

func (a *AmsDos) readDirectories(sectorSize uint16, track *TrackInformation) {
	// 64 files * 32-bytes each = 2048 bytes
	maxDirSectors := (amsdos.DRM * 32) / sectorSize

	// merge the sector data into one slice
	var dirBytes []byte
	for _, s := range track.SectorData[0 : maxDirSectors-1] {
		for _, b := range s {
			dirBytes = append(dirBytes, b)
		}
	}

	// Unmarshal the directory entries
	reader := bytes.NewReader(dirBytes)
	for {
		dir := amsdos.Directory{}
		err := binary.Read(reader, binary.LittleEndian, &dir)
		if err != nil && err == io.EOF {
			break
		} else if err != nil {
			panic("sector read error: " + err.Error())
		}
		a.Directories = append(a.Directories, dir)
	}
}

// Constructs an AMSDOS Extended Disk Parameter Block
func (a *AmsDos) generateDPB(trackSize, sectorSize uint16, firstSectorID, mediaType uint8) {
	dpb := amsdos.DiskParameterBlock{
		ExtentMask:           amsdos.ExtentMask,
		BlockCount:           amsdos.DSM - 1,
		DirectoryCount:       amsdos.DRM - 1,
		Checksum:             0, // CKS = 0 (Fixed Media)
		ReservedTracksOffset: 0, // TODO:michael

		// AMSDOS extended parameters
		MediaType:           mediaType,
		TrackCountPerSide:   amsdos.TrackCount,
		SectorCountPerTrack: amsdos.SectorsPerTrack,
		FirstSectorNumber:   firstSectorID,
		SectorSize:          sectorSize,
		ReadWriteGap:        amsdos.ReadWriteGap,
		FormatGap:           amsdos.FormatGap,
		MultiTrackFlags:     0, // Non multi-track disk
		FreezeFlag:          1, // Non-zero value: use current format
	}

	dpb.RecordsPerTrack = (trackSize - sectorDataStartAddress) / amsdos.CpmRecordSize

	// BLS, BSH, BLM for the Amstrad CPC standard
	blsTable := amsdos.BlsTable[amsdos.BLS]
	dpb.BlockShift = blsTable.BSH
	dpb.BlockMask = blsTable.BLM

	dirsPerBlock := amsdos.BlsTable[amsdos.BLS].Dirs
	reservedBlocks := len(a.Directories)/int(dirsPerBlock) + 1
	dpb.SetAllocationBitmap(reservedBlocks)

	if physicalRecord, ok := amsdos.PhysicalShiftMaskTable[sectorSize]; ok {
		dpb.PhysicalShift = physicalRecord.PSH
		dpb.PhysicalMask = physicalRecord.PHM
	}

	a.DPB = dpb
}
