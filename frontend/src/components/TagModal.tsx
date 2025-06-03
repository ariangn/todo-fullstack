// src/components/TagModal.tsx
import { useEffect, useState } from "react";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogFooter,
} from "@/components/ui/dialog";
import { Label } from "@/components/ui/label";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { fetchTags } from "../services/tagService";

interface TagModalProps {
  mode: "create" | "edit";
  tagId?: string;
  onSave: (name: string) => Promise<void>;
  onDelete?: () => Promise<void>;
  onClose: () => void;
}

export default function TagModal({
  mode,
  tagId,
  onSave,
  onDelete,
  onClose,
}: TagModalProps) {
  const [name, setName] = useState<string>("");
  const [loading, setLoading] = useState<boolean>(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (mode === "edit" && tagId) {
      (async () => {
        const existingTags = await fetchTags();
        const tag = existingTags.find((t) => t.id === tagId);
        if (tag) {
          setName(tag.name);
        }
      })();
    }
  }, [mode, tagId]);

  const handleSave = async () => {
    if (!name.trim()) {
      setError("Name is required");
      return;
    }
    setError(null);
    setLoading(true);
    try {
      await onSave(name.trim());
      onClose();
    } catch (err: unknown) {
      const msg = err instanceof Error ? err.message : String(err);
      setError(msg);
    } finally {
      setLoading(false);
    }
  };

  const handleDelete = async () => {
    if (onDelete) {
      setLoading(true);
      try {
        await onDelete();
        onClose();
      } catch (err: unknown) {
        const msg = err instanceof Error ? err.message : String(err);
        setError(msg);
      } finally {
        setLoading(false);
      }
    }
  };

  return (
    <Dialog open onOpenChange={onClose}>
      <DialogContent className="sm:max-w-md">
        <DialogHeader>
          <DialogTitle>
            {mode === "create" ? "New Tag" : "Edit Tag"}
          </DialogTitle>
        </DialogHeader>

        <div className="space-y-4 mt-2">
          <div className="grid grid-cols-3 items-center gap-4">
            <Label htmlFor="tag-name">Name</Label>
            <Input
              id="tag-name"
              value={name}
              onChange={(e) => setName(e.target.value)}
              className="col-span-2 h-8"
            />
          </div>
          {error && <p className="text-red-500 text-sm">{error}</p>}
        </div>

        <DialogFooter className="space-x-2">
          <Button variant="outline" onClick={onClose} disabled={loading}>
            Cancel
          </Button>
          {mode === "edit" && onDelete && (
            <Button
              variant="destructive"
              onClick={handleDelete}
              disabled={loading}
            >
              Delete
            </Button>
          )}
          <Button onClick={handleSave} disabled={loading}>
            {loading ? "Saving…" : "Save"}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
