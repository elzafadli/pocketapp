package seed

type SeedStatus string

const (
	SEED_STATUS_RUNNING SeedStatus = "running"
	SEED_STATUS_SUCCESS SeedStatus = "success"
	SEED_STATUS_FAILED  SeedStatus = "failed"
)

func (s SeedStatus) String() string {
	return string(s)
}

type SeedTenantType string

const (
	SEED_TENANT_TYPE_DAGANG                     SeedTenantType = "1"
	SEED_TENANT_TYPE_JASA                       SeedTenantType = "2"
	SEED_TENANT_TYPE_KULINER                    SeedTenantType = "3"
	SEED_TENANT_TYPE_DISTRIBUTOR                SeedTenantType = "4"
	SEED_TENANT_TYPE_KOPERASI                   SeedTenantType = "5" // koperasi simpan pinjam
	SEED_TENANT_TYPE_MANUFAKTUR                 SeedTenantType = "6"
	SEED_TENANT_TYPE_MULTI_USAHA                SeedTenantType = "7"
	SEED_TENANT_TYPE_LAINNYA                    SeedTenantType = "8"
	SEED_TENANT_TYPE_KOPERASI_SEKTOR_RIIL       SeedTenantType = "9" // koperasi multi usaha
	SEED_TENANT_TYPE_KOPERASI_PENGADAAN         SeedTenantType = "10"
	SEED_TENANT_TYPE_KOPERASI_KLINIK            SeedTenantType = "11"
	SEED_TENANT_TYPE_KOPERASI_APOTEK            SeedTenantType = "12"
	SEED_TENANT_TYPE_KOPERASI_PERGUDANGAN       SeedTenantType = "13"
	SEED_TENANT_TYPE_KOPERASI_KANTOR            SeedTenantType = "14"
	SEED_TENANT_TYPE_KOPERASI_LOGISTIK          SeedTenantType = "15"
	SEED_TENANT_TYPE_YAYASAN_PENDIDIKAN_KAMPUS  SeedTenantType = "16" // yayasan
	SEED_TENANT_TYPE_YAYASAN_PENDIDIKAN_SEKOLAH SeedTenantType = "17"
	SEED_TENANT_TYPE_KOPERASI_DESA_MERAH_PUTIH  SeedTenantType = "18" // gantari
)

type SeedType string

func (s SeedType) String() string {
	return string(s)
}

const (
	SEED_TYPE_SEEDING SeedType = "seed"
	SEED_TYPE_DEMO    SeedType = "seed_demo"
)

func IsTenantCooperative(tenantType SeedTenantType) bool {
	return tenantType == SEED_TENANT_TYPE_KOPERASI_SEKTOR_RIIL ||
		tenantType == SEED_TENANT_TYPE_KOPERASI_PENGADAAN ||
		tenantType == SEED_TENANT_TYPE_KOPERASI_KLINIK ||
		tenantType == SEED_TENANT_TYPE_KOPERASI_APOTEK ||
		tenantType == SEED_TENANT_TYPE_KOPERASI_PERGUDANGAN ||
		tenantType == SEED_TENANT_TYPE_KOPERASI_KANTOR ||
		tenantType == SEED_TENANT_TYPE_KOPERASI_LOGISTIK ||
		tenantType == SEED_TENANT_TYPE_KOPERASI ||
		tenantType == SEED_TENANT_TYPE_KOPERASI_DESA_MERAH_PUTIH
}

func IsTenantKDMP(tenantType SeedTenantType) bool {
	return tenantType == SEED_TENANT_TYPE_KOPERASI_DESA_MERAH_PUTIH
}

func IsTenantYayasan(tenantType SeedTenantType) bool {
	return tenantType == SEED_TENANT_TYPE_YAYASAN_PENDIDIKAN_KAMPUS ||
		tenantType == SEED_TENANT_TYPE_YAYASAN_PENDIDIKAN_SEKOLAH
}

func IsTenantGeneral(tenantType SeedTenantType) bool {
	return tenantType == SEED_TENANT_TYPE_DAGANG ||
		tenantType == SEED_TENANT_TYPE_JASA ||
		tenantType == SEED_TENANT_TYPE_KULINER ||
		tenantType == SEED_TENANT_TYPE_DISTRIBUTOR ||
		tenantType == SEED_TENANT_TYPE_MANUFAKTUR ||
		tenantType == SEED_TENANT_TYPE_MULTI_USAHA ||
		tenantType == SEED_TENANT_TYPE_LAINNYA
}
