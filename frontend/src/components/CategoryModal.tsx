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
// import type { Category } from "../types";
import { fetchCategories } from "../services/categoryService";

interface CategoryModalProps {
  mode: "create" | "edit";
  categoryId?: string;
  onSave: (name: string, color: string) => Promise<void>;
  onClose: () => void;
}

export default function CategoryModal({
  mode,
  categoryId,
  onSave,
  onClose,
}: CategoryModalProps) {
  const [name, setName] = useState<string>("");
  const [color, setColor] = useState<string>("#D69E2E");
  const [loading, setLoading] = useState<boolean>(false);
  const [error, setError] = useState<string | null>(null);

  // If editing, fetch existing category data
  useEffect(() => {
    if (mode === "edit" && categoryId) {
      (async () => {
        const cats = await fetchCategories();
        const cat = cats.find((c) => c.id === categoryId);
        if (cat) {
          setName(cat.name);
          setColor(cat.color);
        }
      })();
    }
  }, [mode, categoryId]);

  const handleSave = async () => {
    if (!name.trim()) {
      setError("Name is required");
      return;
    }
    setError(null);
    setLoading(true);
    try {
      await onSave(name.trim(), color);
      onClose();
    } catch (err: unknown) {
      const msg = err instanceof Error ? err.message : String(err);
      setError(msg);
    } finally {
      setLoading(false);
    }
  };

  return (
    <Dialog open onOpenChange={onClose}>
      <DialogContent className="sm:max-w-md bg-white">
        <DialogHeader>
          <DialogTitle>
            {mode === "create" ? "New Category" : "Edit Category"}
          </DialogTitle>
        </DialogHeader>

        <div className="space-y-4 mt-2">
          <div className="grid grid-cols-3 items-center gap-4">
            <Label htmlFor="cat-name">Name</Label>
            <Input
              id="cat-name"
              value={name}
              onChange={(e) => setName(e.target.value)}
              className="col-span-2 h-8"
            />
          </div>
          <div className="grid grid-cols-3 items-center gap-4">
            <Label htmlFor="cat-color">Color</Label>
            <Input
              id="cat-color"
              type="color"
              value={color}
              onChange={(e) => setColor(e.target.value)}
              className="col-span-2 h-8 p-0"
            />
          </div>
          {error && <p className="text-red-500 text-sm">{error}</p>}
        </div>

        <DialogFooter>
          <Button
            variant="outline"
            onClick={onClose}
            disabled={loading}
          >
            Cancel
          </Button>
          <Button onClick={handleSave} disabled={loading}>
            {loading ? "Saving…" : "Save"}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
