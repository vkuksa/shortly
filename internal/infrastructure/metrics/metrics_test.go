package metrics_test

// func TestMetricsEndpoint(t *testing.T) {
// 	s := MustOpenMetricServer()
// 	defer MustCloseServer(t, s)

// 	s := MustOpenMetricServer()
// 	defer MustCloseServer(t, s)

// 	testURL := "http://example.com"
// 	testUUID := "test"
// 	exmpl := &domain.Link{
// 		URL:  testURL,
// 		UUID: testUUID,
// 	}
// 	ps.LinkService.GetOriginalLinkFn = func(ctx context.Context, uuid string) (*domain.Link, error) {
// 		if uuid != testUUID {
// 			return nil, usecase.NewError(usecase.ErrNotFound, "Expected")
// 		}

// 		return exmpl, nil
// 	}
// 	ps.LinkService.AddHitFn = func(ctx context.Context, uuid string) error {
// 		return nil
// 	}
// 	ps.LinkService.GenerateShortenedLinkFn = func(ctx context.Context, url string) (*domain.Link, error) {
// 		if url != testURL {
// 			return nil, usecase.NewError(usecase.ErrInvalid, "Expected")
// 		}

// 		return exmpl, nil
// 	}

// 	t.Run("Handle Metrics", func(t *testing.T) {
// 		// Issue request for see if metrics endpoint is accessible
// 		resp, err := http.DefaultClient.Do(ps.MustNewMetricsRequest(t, context.Background(), "GET", nil))
// 		if err != nil {
// 			t.Fatal(err)
// 		} else if got, want := resp.StatusCode, http.StatusOK; got != want {
// 			t.Fatalf("StatusCode=%v, want %v", got, want)
// 		}
// 	})

// 	t.Run("Handle metric requestCount", func(t *testing.T) {
// 		// Make count increment
// 		resp, err := http.DefaultClient.Do(ps.MustNewRequest(t, context.Background(), "GET", "/", nil))
// 		if err != nil {
// 			t.Fatal(err)
// 		} else if got, want := resp.StatusCode, http.StatusOK; got != want {
// 			t.Fatalf("StatusCode=%v, want %v", got, want)
// 		}

// 		// Issue request for querying expected error
// 		resp, err = http.DefaultClient.Do(ps.MustNewMetricsRequest(t, context.Background(), "GET", nil))
// 		if err != nil {
// 			t.Fatal(err)
// 		} else if got, want := resp.StatusCode, http.StatusOK; got != want {
// 			t.Fatalf("StatusCode=%v, want %v", got, want)
// 		}

// 		body, err := io.ReadAll(resp.Body)
// 		if err != nil {
// 			t.Fatal(err)
// 		}

// 		if !strings.Contains(string(body), "http_request_count") {
// 			t.Fatal("Response does not contain required metric")
// 		}
// 	})
// }
