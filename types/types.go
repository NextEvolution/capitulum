package types

import "nextevolution/data-service/types"

type FbLoginReq struct {
	FbToken string `json:"fb_token"`
}

type Config struct {
	Port int `json:"port"`
	StaticFilePath string `json:"static_file_path"`
	Origins []string `json:"origins"`
}

type SalesList struct {
	Sales []Sale `json:"sales"`
	ScanDate int `json:"scan_date"`
}

type Sale struct {
	Customer Customer `json:"customer"`
	Product Product `json:"product"`
	SalesComment SalesComment `json:"sales_comment"`
	Date int
}

type Customer struct {
	Name string `json:"name"`
	FbId string `json:"fb_id"`
}

type Product struct {
	FbId string `json:"fb_id"`
	Album string `json:"album"`
	Description string `json:"description"`
	PictureUrl string `json:"picture_url"`
}

type SalesComment struct {
	Text string `json:"text"`
	FbId string `json:"fb_id"`
}

//TODO: refactor data structure so conversions are not needed
func ConvertDataSales(sas types.SellerAlbumScan) SalesList {
	salesList := SalesList{}
	salesList.ScanDate = sas.Date

	for _, product := range sas.Products{
		for _, salesEvent := range product.SaleEvents {
			sale := &Sale{}
			sale.Date = salesEvent.Date

			sale.Customer.FbId = salesEvent.Customer.Metadata.FbId
			sale.Customer.Name = salesEvent.Customer.Name

			sale.Product.Album = product.Album
			sale.Product.Description = product.Description
			sale.Product.FbId = product.Metadata.FbId
			sale.Product.PictureUrl = product.Metadata.ImageUrl

			sale.SalesComment.FbId = salesEvent.Metadata.FbId
			sale.SalesComment.Text = salesEvent.Metadata.Text
			salesList.Sales = append(salesList.Sales, *sale)
		}
	}
	return salesList
}