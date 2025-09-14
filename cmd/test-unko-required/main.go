package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func main() {
	// テスト1: unko_noなしでリクエスト（エラーになるべき）
	fmt.Println("=== Test 1: Request without unko_no (should fail) ===")
	resp, err := http.Get("http://localhost:8080/dtako/events?from=2025-09-13&to=2025-09-13")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("Status: %d\n", resp.StatusCode)
	fmt.Printf("Response: %s\n", string(body))

	if resp.StatusCode != 400 {
		fmt.Println("❌ ERROR: Expected 400 Bad Request, got", resp.StatusCode)
		fmt.Println("unko_no必須化が正しく動作していません！")
	} else {
		fmt.Println("✅ SUCCESS: unko_no必須チェックが正常に動作しています")
	}

	fmt.Println("\n=== Test 2: Request with unko_no (should succeed) ===")
	resp2, err := http.Get("http://localhost:8080/dtako/events?from=2025-09-13&to=2025-09-13&unko_no=25091310254800000040311")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	defer resp2.Body.Close()

	fmt.Printf("Status: %d\n", resp2.StatusCode)
	if resp2.StatusCode == 200 {
		fmt.Println("✅ SUCCESS: unko_no付きリクエストは正常に処理されました")
	} else {
		body2, _ := io.ReadAll(resp2.Body)
		fmt.Printf("Response: %s\n", string(body2))
	}
}