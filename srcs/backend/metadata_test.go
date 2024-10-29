package backend_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"os"
	"encoding/json"
	"github.com/rsasada/sqluid/srcs/backend"
)

var _ = Describe("MemoryBackend Metadata", func() {
	var memory *backend.MemoryBackend
	var table *backend.Table

	BeforeEach(func() {
		// テスト用のMemoryBackendとTableを初期化
		memory = &backend.MemoryBackend{
			Tables: make(map[string]*Table),
		}
		table = &backend.Table{
			Columns:     []string{"id", "name"},
			ColumnTypes: []backend.ColumnType{backend.IntType, backend.TextType},
			ColumnSize:  []uint{4, 50},
			NumRows:      10,
		}
		memory.Tables["test_table"] = table
	})

	AfterEach(func() {
		// テスト後、ファイルを削除
		os.Remove("TableMeta.json")
	})

	Describe("SaveMetadata", func() {
		It("should save metadata to a JSON file", func() {
			err := memory.SaveMetadata()
			Expect(err).ToNot(HaveOccurred())

			// ファイルが存在するか確認
			_, err = os.Stat("TableMeta.json")
			Expect(err).ToNot(HaveOccurred())

			// ファイル内容を検証
			bytes, err := os.ReadFile("TableMeta.json")
			Expect(err).ToNot(HaveOccurred())

			var metadata backend.Metadata
			err = json.Unmarshal(bytes, &metadata)
			Expect(err).ToNot(HaveOccurred())

			Expect(metadata.Tables).To(HaveLen(1))
			Expect(metadata.Tables[0].Name).To(Equal("test_table"))
			Expect(metadata.Tables[0].Columns).To(Equal([]string{"id", "name"}))
			Expect(metadata.Tables[0].ColumnTypes).To(Equal([]backend.ColumnType{IntType, TextType}))
			Expect(metadata.Tables[0].ColumnSize).To(Equal([]uint{4, 50}))
			Expect(metadata.Tables[0].NumRows).To(Equal(uint(10)))
		})
	})

	Describe("LoadMetadata", func() {
		BeforeEach(func() {
			// SaveMetadataで事前にメタデータをファイルに保存
			err := memory.SaveMetadata()
			Expect(err).ToNot(HaveOccurred())
		})

		It("should load metadata from a JSON file", func() {
			// 既存のテーブルデータをクリア
			memory.Tables = make(map[string]*backend.Table)

			err := memory.LoadMetadata()
			Expect(err).ToNot(HaveOccurred())

			// メタデータが正しくロードされているか確認
			Expect(memory.Tables).To(HaveKey("test_table"))
			loadedTable := memory.Tables["test_table"]
			Expect(loadedTable.Columns).To(Equal([]string{"id", "name"}))
			Expect(loadedTable.ColumnTypes).To(Equal([]backend.ColumnType{backend.IntType, backend.TextType}))
			Expect(loadedTable.ColumnSize).To(Equal([]uint{4, 50}))
			Expect(loadedTable.RowNum).To(Equal(uint(10)))
		})
	})
})
