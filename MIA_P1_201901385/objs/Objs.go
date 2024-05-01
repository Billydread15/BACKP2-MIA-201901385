package objs

import (
	"fmt"
)

type MBR struct {
	Mbr_tamano         int32
	Mbr_fecha_creacion [10]byte
	Mbr_dsk_signature  int32
	Mbr_dsk_fit        [1]byte
	Mbr_partitions     [4]Partition
}

func PrintMBR(mbr MBR) {
	fmt.Println(" ")
	fmt.Println("mbr_tamano:", mbr.Mbr_tamano)
	fmt.Println("mbr_fecha_creacion:", string(mbr.Mbr_fecha_creacion[:]))
	fmt.Println("mbr_dsk_signature:", mbr.Mbr_dsk_signature)
	fmt.Println("mbr_dsk_fit:", string(mbr.Mbr_dsk_fit[:]))
	for i := 0; i < 4; i++ {
		PrintPartition(mbr.Mbr_partitions[i])
	}
}

func ListPartitions(mbr MBR) []string {
	var partitions []string
	for i := 0; i < 4; i++ {
		partitions = append(partitions, ReturnPartitionName(mbr.Mbr_partitions[i]))
	}
	return partitions
}

type Partition struct {
	Part_status      [1]byte
	Part_type        [1]byte
	Part_fit         [1]byte
	Part_start       int32
	Part_s           int32
	Part_name        [16]byte
	Part_correlative int32
	Part_id          [4]byte
}

func PrintPartition(part Partition) {
	fmt.Print("status:", string(part.Part_status[:]))
	fmt.Print(" type:", string(part.Part_type[:]))
	fmt.Print(" fit:", string(part.Part_fit[:]))
	fmt.Print(" start:", part.Part_start)
	fmt.Print(" s:", part.Part_s)
	fmt.Print(" name:", string(part.Part_name[:]))
	fmt.Print(" correlative:", part.Part_correlative)
	fmt.Println(" id:", string(part.Part_id[:]))
}

func ReturnPartitionName(part Partition) string {
	return string(part.Part_name[:])
}

type EBR struct {
	Part_mount [1]byte
	Part_fit   [1]byte
	Part_start int32
	Part_s     int32
	Part_next  int32
	Part_name  [16]byte
}

func PrintEBR(ebr EBR) {
	fmt.Print("mount:", string(ebr.Part_mount[:]))
	fmt.Print(" fit:", string(ebr.Part_fit[:]))
	fmt.Print(" start:", ebr.Part_start)
	fmt.Print(" size:", ebr.Part_s)
	fmt.Print(" next:", ebr.Part_next)
	fmt.Println(" name:", string(ebr.Part_name[:]))
}

type Superblock struct {
	S_filesystem_type   int32
	S_inodes_count      int32
	S_blocks_count      int32
	S_free_blocks_count int32
	S_free_inodes_count int32
	S_mtime             [16]byte
	S_umtime            [16]byte
	S_mnt_count         int32
	S_magic             int32
	S_inode_s           int32
	S_block_s           int32
	S_fist_ino          int32
	S_first_blo         int32
	S_bm_inode_start    int32
	S_bm_block_start    int32
	S_inode_start       int32
	S_block_start       int32
}

type Inode struct {
	I_uid   int32
	I_gid   int32
	I_s     int32
	I_atime [16]byte
	I_ctime [16]byte
	I_mtime [16]byte
	I_block [15]int32
	I_type  [1]byte
	I_perm  [3]byte
}

type Folderblock struct {
	B_content [4]Content
}

type Content struct {
	B_name  [12]byte
	B_inodo int32
}

type Fileblock struct {
	B_content [64]byte
}

type Pointerblock struct {
	B_pointers [16]int32
}

type Journaling struct {
	J_size      int32
	J_ultimo    int32
	J_contenido [50]Content_j
}

type Content_j struct {
	Operation [10]byte
	Path      [100]byte
	Content   [100]byte
	Date      [16]byte
}
