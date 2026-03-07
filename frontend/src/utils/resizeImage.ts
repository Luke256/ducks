/**
 * クライアントサイドで画像をリサイズしてからアップロード用のFileを返すユーティリティ。
 * Canvas APIを使用して、指定した最大幅・最大高さに収まるようリサイズし、
 * WebP形式（非対応ブラウザではJPEG）で圧縮する。
 */

const DEFAULT_MAX_WIDTH = 800;
const DEFAULT_MAX_HEIGHT = 800;
const DEFAULT_QUALITY = 0.8;

export interface ResizeOptions {
  /** リサイズ後の最大幅 (デフォルト: 1920) */
  maxWidth?: number;
  /** リサイズ後の最大高さ (デフォルト: 1920) */
  maxHeight?: number;
  /** 圧縮品質 0.0〜1.0 (デフォルト: 0.8) */
  quality?: number;
}

/**
 * 画像ファイルをリサイズ・圧縮して新しいFileオブジェクトを返す。
 * 元画像が最大サイズ以下の場合でも圧縮は行われる。
 */
export async function resizeImage(
  file: File,
  options: ResizeOptions = {}
): Promise<File> {
  const {
    maxWidth = DEFAULT_MAX_WIDTH,
    maxHeight = DEFAULT_MAX_HEIGHT,
    quality = DEFAULT_QUALITY,
  } = options;

  // 画像をImageElementに読み込む
  const imageBitmap = await createImageBitmap(file);
  const { width: origWidth, height: origHeight } = imageBitmap;

  // アスペクト比を維持しつつ最大サイズに収める
  let newWidth = origWidth;
  let newHeight = origHeight;

  if (origWidth > maxWidth || origHeight > maxHeight) {
    const ratio = Math.min(maxWidth / origWidth, maxHeight / origHeight);
    newWidth = Math.round(origWidth * ratio);
    newHeight = Math.round(origHeight * ratio);
  }

  // OffscreenCanvas が使える場合はそちらを使う（Web Worker対応）
  // 使えない場合は通常のCanvasを使う
  const canvas = document.createElement("canvas");
  canvas.width = newWidth;
  canvas.height = newHeight;

  const ctx = canvas.getContext("2d");
  if (!ctx) {
    throw new Error("Canvas 2D context の取得に失敗しました");
  }

  ctx.drawImage(imageBitmap, 0, 0, newWidth, newHeight);
  imageBitmap.close();

  // WebP対応チェック → 非対応ならJPEGにフォールバック
  const mimeType = supportsWebP() ? "image/webp" : "image/jpeg";
  const extension = mimeType === "image/webp" ? ".webp" : ".jpg";

  // Canvasからblobを取得
  const blob = await new Promise<Blob>((resolve, reject) => {
    canvas.toBlob(
      (b) => {
        if (b) resolve(b);
        else reject(new Error("画像の圧縮に失敗しました"));
      },
      mimeType,
      quality
    );
  });

  // 元のファイル名から拡張子を差し替え
  const baseName = file.name.replace(/\.[^.]+$/, "");
  const newFileName = baseName + extension;

  return new File([blob], newFileName, { type: mimeType });
}

/** WebPフォーマットのサポートを確認 */
function supportsWebP(): boolean {
  try {
    const canvas = document.createElement("canvas");
    canvas.width = 1;
    canvas.height = 1;
    return canvas.toDataURL("image/webp").startsWith("data:image/webp");
  } catch {
    return false;
  }
}
