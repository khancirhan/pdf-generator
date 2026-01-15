package pdfgen

import (
	"context"
	"sync"
	"time"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/cdproto/runtime"
	"github.com/chromedp/chromedp"
)

type ChromePDFGenerator struct {
	allocCtx context.Context
	cancel   context.CancelFunc
	mu       sync.Mutex
	ctx      context.Context // Single tab, reused
}

func NewChromePDFGenerator() (*ChromePDFGenerator, error) {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		// Memory optimization flags
		chromedp.Flag("headless", true),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("disable-dev-shm-usage", true),

		// Disable unnecessary features - BIG memory savings
		chromedp.Flag("disable-extensions", true),
		chromedp.Flag("disable-plugins", true),
		chromedp.Flag("disable-background-networking", true),
		chromedp.Flag("disable-sync", true),
		chromedp.Flag("disable-translate", true),
		chromedp.Flag("disable-default-apps", true),
		chromedp.Flag("disable-breakpad", true), // Disable crash reporting
		chromedp.Flag("disable-component-extensions-with-background-pages", true),
		chromedp.Flag("disable-features=TranslateUI", true),
		chromedp.Flag("disable-ipc-flooding-protection", true),
		chromedp.Flag("disable-renderer-backgrounding", true),
		chromedp.Flag("enable-features=NetworkService,NetworkServiceInProcess", true),
		chromedp.Flag("force-color-profile", "srgb"),

		// Memory limits
		chromedp.Flag("js-flags", "--max-old-space-size=128"), // Limit V8 heap to 128MB
		chromedp.Flag("disable-hang-monitor", true),
		chromedp.Flag("disable-backgrounding-occluded-windows", true),

		// Single process mode - less memory, slightly more CPU
		chromedp.Flag("single-process", false), // Keep false for stability, true saves ~50MB but risky

		// Reduce rendering overhead
		chromedp.Flag("disable-smooth-scrolling", true),
		chromedp.Flag("disable-lcd-text", true),
		chromedp.Flag("hide-scrollbars", true),
		chromedp.Flag("mute-audio", true),
	)

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	browserCtx, _ := chromedp.NewContext(allocCtx)

	// Start browser
	if err := chromedp.Run(browserCtx); err != nil {
		cancel()
		return nil, err
	}

	// Create single reusable tab
	tabCtx, _ := chromedp.NewContext(browserCtx)
	chromedp.Run(tabCtx, chromedp.Navigate("about:blank"))

	return &ChromePDFGenerator{
		allocCtx: allocCtx,
		cancel:   cancel,
		ctx:      tabCtx,
	}, nil
}

// JavaScript to wait for all images
const waitForImagesScript = `
new Promise((resolve) => {
    const images = Array.from(document.querySelectorAll('img'));
    if (images.length === 0) {
        resolve(true);
        return;
    }
    
    let loaded = 0;
    const total = images.length;
    
    const checkDone = () => {
        loaded++;
        if (loaded >= total) resolve(true);
    };
    
    images.forEach(img => {
        if (img.complete) {
            checkDone();
        } else {
            img.onload = checkDone;
            img.onerror = checkDone; // Don't hang on broken images
        }
    });
    
    // Safety timeout
    setTimeout(() => resolve(true), 10000);
})
`

// More thorough - waits for images, fonts, and idle network
const waitForFullLoadScript = `
Promise.all([
    // Wait for images
    Promise.all(
        Array.from(document.querySelectorAll('img'))
            .filter(img => !img.complete)
            .map(img => new Promise(resolve => {
                img.onload = resolve;
                img.onerror = resolve;
            }))
    ),
    // Wait for fonts
    document.fonts.ready,
]).then(() => true)
`

func (g *ChromePDFGenerator) GeneratePDF(html string, opts PDFOptions) ([]byte, error) {
	g.mu.Lock()
	defer g.mu.Unlock()

	var pdf []byte

	err := chromedp.Run(g.ctx,
		// Set content
		chromedp.ActionFunc(func(ctx context.Context) error {
			frameTree, err := page.GetFrameTree().Do(ctx)
			if err != nil {
				return err
			}
			return page.SetDocumentContent(frameTree.Frame.ID, html).Do(ctx)
		}),

		// Wait for images and fonts to load
		chromedp.ActionFunc(func(ctx context.Context) error {
			_, exp, err := runtime.Evaluate(waitForFullLoadScript).
				WithAwaitPromise(true).
				Do(ctx)
			if err != nil {
				return err
			}
			if exp != nil {
				return exp
			}
			return nil
		}),

		// Wait for custom expression if provided
		chromedp.ActionFunc(func(ctx context.Context) error {
			if opts.WaitForExpression != "" {
				_, exp, err := runtime.Evaluate(opts.WaitForExpression).
					WithAwaitPromise(true).
					Do(ctx)
				if err != nil {
					return err
				}
				if exp != nil {
					return exp
				}
			}
			return nil
		}),

		// Wait for delay if provided
		chromedp.ActionFunc(func(ctx context.Context) error {
			if opts.WaitDelay != "" {
				d, err := time.ParseDuration(opts.WaitDelay)
				if err != nil {
					return err
				}
				time.Sleep(d)
			}
			return nil
		}),

		// Generate PDF
		chromedp.ActionFunc(func(ctx context.Context) error {
			var err error
			pdf, _, err = page.PrintToPDF().
				WithPrintBackground(opts.PrintBackground).
				WithPreferCSSPageSize(opts.PreferCSSPageSize).
				WithLandscape(opts.Landscape).
				WithPaperWidth(opts.PaperWidth).
				WithPaperHeight(opts.PaperHeight).
				WithMarginTop(opts.MarginTop).
				WithMarginBottom(opts.MarginBottom).
				WithMarginLeft(opts.MarginLeft).
				WithMarginRight(opts.MarginRight).
				Do(ctx)
			return err
		}),

		// Clear page to free memory
		// chromedp.Navigate("about:blank"),
	)

	return pdf, err
}

func (g *ChromePDFGenerator) Close() {
	g.cancel()
}
