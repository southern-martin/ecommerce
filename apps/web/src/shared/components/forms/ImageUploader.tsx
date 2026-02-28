import { useState, useCallback, useRef } from "react";
import { Upload, X, ImageIcon, Loader2 } from "lucide-react";
import { Button } from "@/shared/components/ui/button";
import { Badge } from "@/shared/components/ui/badge";
import { cn } from "@/shared/lib/utils";
import { uploadImage } from "@/shared/services/media.api";

interface ImageUploaderProps {
  value?: string[];
  onChange?: (urls: string[]) => void;
  maxFiles?: number;
  maxSizeMB?: number;
  accept?: string;
  className?: string;
  ownerType?: string;
}

export function ImageUploader({
  value = [],
  onChange,
  maxFiles = 5,
  maxSizeMB = 5,
  accept = "image/*",
  className,
  ownerType = "product",
}: ImageUploaderProps) {
  const [isDragging, setIsDragging] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [uploading, setUploading] = useState<number>(0);
  const inputRef = useRef<HTMLInputElement>(null);

  const handleFiles = useCallback(
    async (fileList: FileList | null) => {
      if (!fileList) return;
      setError(null);

      const files = Array.from(fileList);

      if (value.length + files.length > maxFiles) {
        setError(`Maximum ${maxFiles} images allowed`);
        return;
      }

      for (const file of files) {
        if (file.size > maxSizeMB * 1024 * 1024) {
          setError(`File ${file.name} exceeds ${maxSizeMB}MB limit`);
          return;
        }
      }

      setUploading(files.length);

      const newUrls: string[] = [];
      for (const file of files) {
        try {
          const media = await uploadImage(file, ownerType);
          newUrls.push(media.url);
        } catch {
          setError(`Failed to upload ${file.name}`);
        }
      }

      setUploading(0);

      if (newUrls.length > 0) {
        onChange?.([...value, ...newUrls]);
      }
    },
    [value, maxFiles, maxSizeMB, onChange, ownerType]
  );

  const handleDrop = useCallback(
    (e: React.DragEvent) => {
      e.preventDefault();
      setIsDragging(false);
      handleFiles(e.dataTransfer.files);
    },
    [handleFiles]
  );

  const handleDragOver = useCallback((e: React.DragEvent) => {
    e.preventDefault();
    setIsDragging(true);
  }, []);

  const handleDragLeave = useCallback((e: React.DragEvent) => {
    e.preventDefault();
    setIsDragging(false);
  }, []);

  const removeImage = (url: string) => {
    onChange?.(value.filter((u) => u !== url));
  };

  return (
    <div className={cn("space-y-3", className)}>
      {/* Drop Zone */}
      <div
        onDrop={handleDrop}
        onDragOver={handleDragOver}
        onDragLeave={handleDragLeave}
        onClick={() => inputRef.current?.click()}
        className={cn(
          "flex cursor-pointer flex-col items-center justify-center rounded-lg border-2 border-dashed p-6 text-center transition-colors",
          isDragging
            ? "border-primary bg-primary/5"
            : "border-muted-foreground/25 hover:border-muted-foreground/50",
          uploading > 0 && "pointer-events-none opacity-60"
        )}
      >
        {uploading > 0 ? (
          <>
            <Loader2 className="mb-2 h-8 w-8 animate-spin text-muted-foreground" />
            <p className="text-sm font-medium">
              Uploading {uploading} file{uploading > 1 ? "s" : ""}...
            </p>
          </>
        ) : (
          <>
            <Upload className="mb-2 h-6 w-6 text-muted-foreground" />
            <p className="text-sm font-medium">
              Drag & drop images, or click to select
            </p>
            <p className="mt-1 text-xs text-muted-foreground">
              Max {maxFiles} files, up to {maxSizeMB}MB each
            </p>
          </>
        )}
        <input
          ref={inputRef}
          type="file"
          accept={accept}
          multiple
          className="hidden"
          onChange={(e) => {
            handleFiles(e.target.files);
            e.target.value = "";
          }}
        />
      </div>

      {error && <p className="text-sm text-destructive">{error}</p>}

      {/* Previews */}
      {value.length > 0 && (
        <div className="grid grid-cols-3 gap-3 sm:grid-cols-4">
          {value.map((url, index) => (
            <div key={url} className="group relative aspect-square">
              <img
                src={url}
                alt={`Image ${index + 1}`}
                className="h-full w-full rounded-md object-cover border"
                onError={(e) => {
                  (e.target as HTMLImageElement).style.display = "none";
                }}
              />
              {index === 0 && (
                <Badge className="absolute bottom-1 left-1 text-[10px] px-1.5 py-0" variant="secondary">
                  Primary
                </Badge>
              )}
              <Button
                type="button"
                variant="destructive"
                size="icon"
                className="absolute -right-1.5 -top-1.5 h-5 w-5 opacity-0 group-hover:opacity-100 transition-opacity"
                onClick={(e) => {
                  e.stopPropagation();
                  removeImage(url);
                }}
              >
                <X className="h-3 w-3" />
              </Button>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
