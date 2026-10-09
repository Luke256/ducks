import { afterEach, expect, it, vi } from "vitest";
import { resizeImage } from "@/utils/resizeImage";
afterEach(() => {
    vi.restoreAllMocks();
    vi.unstubAllGlobals();
});
it("縦横比を保って800px以内に圧縮し、画像リソースを解放する", async () => {
    const close = vi.fn();
    const drawImage = vi.fn();
    vi.stubGlobal(
        "createImageBitmap",
        vi.fn().mockResolvedValue({ width: 1600, height: 1200, close }),
    );
    vi.spyOn(HTMLCanvasElement.prototype, "getContext").mockReturnValue({
        drawImage,
    } as unknown as CanvasRenderingContext2D);
    vi.spyOn(HTMLCanvasElement.prototype, "toDataURL").mockReturnValue(
        "data:image/webp;base64,",
    );
    vi.spyOn(HTMLCanvasElement.prototype, "toBlob").mockImplementation(
        (callback) => callback(new Blob(["compressed"], { type: "image/webp" })),
    );
    const result = await resizeImage(
        new File(["original"], "photo.png", { type: "image/png" }),
    );
    expect(result.name).toBe("photo.webp");
    expect(result.type).toBe("image/webp");
    expect(drawImage).toHaveBeenCalledWith(expect.anything(), 0, 0, 800, 600);
    expect(close).toHaveBeenCalledOnce();
});
it("非対応ブラウザではJPEGにし、小さい画像を拡大しない", async () => {
    const drawImage = vi.fn();
    vi.stubGlobal(
        "createImageBitmap",
        vi.fn().mockResolvedValue({ width: 100, height: 200, close: vi.fn() }),
    );
    vi.spyOn(HTMLCanvasElement.prototype, "getContext").mockReturnValue({
        drawImage,
    } as unknown as CanvasRenderingContext2D);
    vi.spyOn(HTMLCanvasElement.prototype, "toDataURL").mockReturnValue(
        "data:image/png;base64,",
    );
    vi.spyOn(HTMLCanvasElement.prototype, "toBlob").mockImplementation(
        (callback) => callback(new Blob(["jpeg"], { type: "image/jpeg" })),
    );
    const result = await resizeImage(
        new File(["image"], "photo.png", { type: "image/png" }),
    );
    expect(result.name).toBe("photo.jpg");
    expect(drawImage).toHaveBeenCalledWith(expect.anything(), 0, 0, 100, 200);
});
it("画像でない入力やCanvas失敗を知らせ、失敗時もリソースを解放する", async () => {
    await expect(
        resizeImage(new File(["text"], "note.txt", { type: "text/plain" })),
    ).rejects.toThrow("画像ファイル");
    const close = vi.fn();
    vi.stubGlobal(
        "createImageBitmap",
        vi.fn().mockResolvedValue({ width: 100, height: 100, close }),
    );
    vi.spyOn(HTMLCanvasElement.prototype, "getContext").mockReturnValue(null);
    await expect(
        resizeImage(new File(["image"], "photo.png", { type: "image/png" })),
    ).rejects.toThrow("Canvas");
    expect(close).toHaveBeenCalledOnce();
});
