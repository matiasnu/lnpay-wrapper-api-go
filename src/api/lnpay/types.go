package lnpay

// WebhookServer types

type Invoice struct {
	Destination string `json:"destination"`
	PaymentHash string `json:"payment_hash"`
	NumSatoshis string `json:"num_satoshis"`
	Timestamp   string `json:"timestamp"`
	Expiry      string `json:"expiry"`
	Description string `json:"description"`
	CltvExpiry  string `json:"cltv_expiry"`
	RouteHints  []struct {
		HopHints []struct {
			NodeID                    string `json:"node_id"`
			ChanID                    string `json:"chan_id"`
			FeeProportionalMillionths int    `json:"fee_proportional_millionths"`
			CltvExpiryDelta           int    `json:"cltv_expiry_delta"`
		} `json:"hop_hints"`
	} `json:"route_hints"`
	PaymentAddr string `json:"payment_addr"`
	NumMsat     string `json:"num_msat"`
	Features    struct {
		Num9 struct {
			Name    string `json:"name"`
			IsKnown bool   `json:"is_known"`
		} `json:"9"`
		Num15 struct {
			Name    string `json:"name"`
			IsKnown bool   `json:"is_known"`
		} `json:"15"`
		Num17 struct {
			Name    string `json:"name"`
			IsKnown bool   `json:"is_known"`
		} `json:"17"`
	} `json:"features"`
}

type QueryRoutes struct {
	Routes []struct {
		Totaltimelock int    `json:"totalTimeLock"`
		Totalamt      string `json:"totalAmt"`
		Hops          []struct {
			Chanid           string `json:"chanId"`
			Chancapacity     string `json:"chanCapacity"`
			Amttoforward     string `json:"amtToForward"`
			Expiry           int    `json:"expiry"`
			Amttoforwardmsat string `json:"amtToForwardMsat"`
			Feemsat          string `json:"feeMsat,omitempty"`
			Pubkey           string `json:"pubKey"`
			Tlvpayload       bool   `json:"tlvPayload"`
		} `json:"hops"`
		Totalfeesmsat string `json:"totalFeesMsat"`
		Totalamtmsat  string `json:"totalAmtMsat"`
	} `json:"routes"`
	Successprob float64 `json:"successProb"`
}

type Webhook struct {
	CreatedAt int    `json:"created_at"`
	ID        string `json:"id"`
	Event     Event  `json:"event"`
}

// "wallet_created" webhook payload
// https://docs.lnpay.co/webhooks/getting-started#payloads
type WebhookWalletCreated struct {
	Webhook
	Data struct {
		Wal Wallet `json:"wal"`
	} `json:"data"`
}

// "wallet_send" webhook payload
// https://docs.lnpay.co/webhooks/getting-started#payloads
type WebhookWalletSend struct {
	Webhook
	Data struct {
		Wtx Wtx `json:"wtx"`
	} `json:"data"`
}

// "wallet_receive" webhook payload
// https://docs.lnpay.co/webhooks/getting-started#payloads
type WebhookWalletInternalTransfer struct {
	Webhook
	Data struct {
		Wtx Wtx `json:"wtx"`
	} `json:"data"`
}

// "wallet_transfer_IN/OUT" webhook payload
// https://docs.lnpay.co/webhooks/getting-started#payloads
type WebhookWalletReceive struct {
	Webhook
	Data struct {
		Wtx Wtx `json:"wtx"`
	} `json:"data"`
}

// "paywall_created" webhook payload
// https://docs.lnpay.co/webhooks/getting-started#paywalls
type WebhookPaywallCreated struct {
	Webhook
	Data struct {
		Pywl Pywl `json:"pywl"`
	} `json:"data"`
}

// "paywall_conversion" webhook payload
// https://docs.lnpay.co/webhooks/getting-started#paywalls
type WebhookPaywallConversion struct {
	Webhook
	Data struct {
		Pywl Pywl `json:"pywl"`
		Wtx  Wtx  `json:"wtx"`
	} `json:"data"`
}

// Primary types

type Error struct {
	Name    string `json:"name"`
	Message string `json:"message"`
	Code    int    `json:"code"`
	Status  int    `json:"status"`
}

func (err Error) Error() string {
	return err.Message
}

type Wallet struct {
	*Client      `gorm:"-"`
	BaseUrl      string
	ID           string     `json:"id" gorm:"id"`
	UserLabel    string     `json:"user_label" gorm:"user_label"`
	CreatedAt    int        `json:"created_at" gorm:"created_at"`
	UpdatedAt    int        `json:"updated_at" gorm:"updated_at"`
	Balance      int64      `json:"balance" gorm:"balance"`
	StatusType   StatusType `json:"statusType" gorm:"statusType"`
	AccessKeys   AccessKeys `json:"access_keys" gorm:"access_keys"`
	StatusTypeId int
	AccessKeysId int
}

type Wtx struct {
	UserLabel   string                 `json:"user_label"`
	CreatedAt   int                    `json:"created_at"`
	NumSatoshis int64                  `json:"num_satoshis"`
	ID          string                 `json:"id"`
	Wal         Wallet                 `json:"wal"`
	WtxType     WtxType                `json:"wtxType"`
	LnTx        LnTx                   `json:"lnTx"`
	PassThru    map[string]interface{} `json:"passThru"`
}

type LnTx struct {
	ID              string                 `json:"id"`
	CreatedAt       int                    `json:"created_at"`
	DestPubkey      string                 `json:"dest_pubkey"`
	PaymentRequest  string                 `json:"payment_request"`
	RHashDecoded    string                 `json:"r_hash_decoded"`
	Memo            string                 `json:"memo"`
	DescriptionHash string                 `json:"description_hash"`
	NumSatoshis     int64                  `json:"num_satoshis"`
	Expiry          int                    `json:"expiry"`
	ExpiresAt       int                    `json:"expires_at"`
	PaymentPreimage string                 `json:"payment_preimage"`
	Settled         int                    `json:"settled"`
	SettledAt       int                    `json:"settled_at"`
	IsKeysend       bool                   `json:"is_keysend"`
	CustomRecords   map[string]interface{} `json:"custom_records"`
}

type Pywl struct {
	DestinationURL string                 `json:"destination_url"`
	Memo           string                 `json:"memo"`
	ShortURL       string                 `json:"short_url"`
	NumSatoshis    int64                  `json:"lnd_value"`
	CreatedAt      int                    `json:"created_at"`
	UpdatedAt      int                    `json:"updated_at"`
	Metadata       map[string]interface{} `json:"metadata"`
	ID             string                 `json:"id"`
	PaywallLink    string                 `json:"paywall_link"`
	CustyDomain    struct {
		DomainName string `json:"domain_name"`
	} `json:"custyDomain"`
	StatusType  StatusType `json:"statusType"`
	PaywallType struct {
		Name        string `json:"name"`
		DisplayName string `json:"display_name"`
		Description string `json:"description"`
	} `json:"paywallType"`
	Template struct {
		Layout string `json:"layout"`
	} `json:"template"`
	LinkExpRule struct {
		Type        string `json:"type"`
		Name        string `json:"name"`
		DisplayName string `json:"display_name"`
		TimeMinutes int    `json:"time_minutes"`
	} `json:"linkExpRule"`
}

// Secondary/helper types

type Event struct {
	Type        string `json:"type"`
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
}

type AccessKeys struct {
	WalletAdmin   []string `json:"Wallet Admin" gorm:"type:text[]"`
	WalletInvoice []string `json:"Wallet Invoice" gorm:"type:text[]"`
	WalletRead    []string `json:"Wallet Read" gorm:"type:text[]"`
}

type StatusType struct {
	Type        string `json:"type"`
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
}

type WtxType struct {
	Layer       string `json:"layer"`
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
}
