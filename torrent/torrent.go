// https://en.wikipedia.org/wiki/Torrent_file

package torrent

// A torrent file contains a list of files and integrity metadata about all the pieces, and optionally contains a large list of trackers.
// A torrent file is a bencoded dictionary with the following keys (the keys in any bencoded dictionary are lexicographically ordered):
type TorrentFile struct {
	// the URL of the high tracker
	Announce string
	// this maps to a dictionary whose keys are very dependent on whether one or more files are being shared:
	Info
}

type Info struct {
	// a list of dictionaries each corresponding to a file (only when multiple files are being shared). Each dictionary has the following keys:
	Files []File
	// size of the file in bytes (only when one file is being shared though)
	Length int
	// suggested filename where the file is to be saved (if one file)/suggested directory name where the files are to be saved (if multiple files)
	Name string
	// number of bytes per piece. This is commonly 28 KiB = 256 KiB = 262,144 B.
	PieceLength int
	// a hash list, i.e., a concatenation of each piece's SHA-1 hash. As SHA-1 returns a 160-bit hash, pieces will be a string whose length is a multiple of 20 bytes. If the torrent contains multiple files, the pieces are formed by concatenating the files in the order they appear in the files dictionary (i.e., all pieces in the torrent are the full piece length except for the last piece, which may be shorter).
	Pieces []byte
}

type File struct {
	// size of the file in bytes
	Length int
	// a list of strings corresponding to subdirectory names, the last of which is the actual file name
	Path []string
}
