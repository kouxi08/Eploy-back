package cloudflare

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/cloudflare/cloudflare-go"
	"github.com/kouxi08/Eploy/utils"
)

type CloudflareAPP struct {
	app *cloudflare.API
}

func InitCloudflareApp() (*CloudflareAPP, error) {
	app, err := cloudflare.New(os.Getenv("CLOUDFLARE_API_KEY"), os.Getenv("CLOUDFLARE_API_EMAIL"))
	if err != nil {
		log.Fatalf("error initializing Cloudflare app: %v\n", err)
	}
	return &CloudflareAPP{app}, nil
}

func (c *CloudflareAPP) listRecord(zoneID string) {
	records, _, err := c.app.ListDNSRecords(context.Background(), cloudflare.ZoneIdentifier(zoneID), cloudflare.ListDNSRecordsParams{})
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, r := range records {
		fmt.Printf("%s: %s: %s\n", r.ID, r.Name, r.Content)
	}
}

func (c *CloudflareAPP) AddRecord(record *utils.DNSRecord, name string) error {

	zoneID, err := c.app.ZoneIDByName(record.Domain)
	if err != nil {
		return err
	}

	proxied := record.Proxied
	dnsrecord := cloudflare.CreateDNSRecordParams{
		Type:    record.Type,
		Name:    name,
		Content: record.Content,
		TTL:     record.TTL,
		Proxied: &proxied,
	}

	_, err = c.app.CreateDNSRecord(context.Background(), cloudflare.ZoneIdentifier(zoneID), dnsrecord)
	if err != nil {
		return err
	}
	return nil
}

func (c *CloudflareAPP) DeleteRecord(record *utils.DNSRecord, name string) error {

	zoneID, err := c.app.ZoneIDByName(record.Domain)
	if err != nil {
		return err
	}
	records, _, err := c.app.ListDNSRecords(context.Background(), cloudflare.ZoneIdentifier(zoneID), cloudflare.ListDNSRecordsParams{
		Name: fmt.Sprintf("%s.%s", name, record.Domain),
		Type: record.Type,
	})
	if err != nil {
		log.Fatal(err)
	}
	dnsRecordID := records[0].ID

	err = c.app.DeleteDNSRecord(context.Background(), cloudflare.ZoneIdentifier(zoneID), dnsRecordID)
	if err != nil {
		return err
	}

	return nil
}
