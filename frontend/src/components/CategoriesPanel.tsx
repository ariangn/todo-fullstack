import { useState } from "react";
import { Popover, PopoverTrigger, PopoverContent } from "@/components/ui/popover";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { PlusIcon, XIcon, Pencil } from "lucide-react";
import { createCategory } from "../services/categoryService";
import type { Category } from "../types";

interface CategoriesPanelProps {
  categories: Category[];
  onEdit: (id: string) => void;
  onDelete: (id: string) => void;
  refreshCategories: () => void;
}

export default function CategoriesPanel({
  categories,
  onEdit,
  onDelete,
  refreshCategories,
}: CategoriesPanelProps) {
  const [isPopoverOpen, setIsPopoverOpen] = useState(false);
  const [newName, setNewName] = useState<string>("");
  const [newColor, setNewColor] = useState<string>("#D8D8F6");
  const [loading, setLoading] = useState<boolean>(false);
  const [error, setError] = useState<string | null>(null);

  const handleCreate = async () => {
    if (!newName.trim()) {
      setError("Name is required");
      return;
    }
    setError(null);
    setLoading(true);

    try {
      await createCategory(newName.trim(), newColor);
      setNewName("");
      setNewColor("#D8D8F6");
      refreshCategories();
      // Close popover by blurring the trigger button
      document.getElementById("add-cat-trigger")?.blur();
    } catch (err: unknown) {
      // If err is an Error, extract message; otherwise string
      const msg = err instanceof Error ? err.message : String(err);
      setError(msg);
    } finally {
      setLoading(false);
    }
  };

  return (
    <Card className="w-full max-w-sm flex flex-col self-start">
      <CardHeader>
        <CardTitle>Categories</CardTitle>
        <CardDescription>
          Create, edit, or delete categories to organize your items.
        </CardDescription>
      </CardHeader>

      <CardContent className="flex-1 overflow-y-auto bg-white">
        <ul className="space-y-2 mb-4">
          {(categories ?? []).map((cat) => (
            <li
              key={cat.id}
              className="flex items-center p-2 rounded hover:bg-gray-100 cursor-pointer"
            >
              <span
                className="inline-block h-4 w-4 rounded-full mr-2"
                style={{ backgroundColor: cat.color }}
              />
              <span className="flex-1 text-sm text-gray-900">{cat.name}</span>
              <div className="flex space-x-1">
                <Button variant="ghost" size="icon" onClick={() => onEdit(cat.id)}>
                  <Pencil className="h-4 w-4 text-gray-600 hover:text-gray-800" />
                </Button>
                <Button variant="ghost" size="icon" onClick={() => onDelete(cat.id)}>
                  <XIcon className="h-4 w-4 text-red-600 hover:text-red-800" />
                </Button>
              </div>
            </li>
          ))}
        </ul>

        <div className="flex justify-center w-full">
          <Popover open={isPopoverOpen} onOpenChange={setIsPopoverOpen}>
            <PopoverTrigger asChild>
              <Button
                id="add-cat-trigger"
                className="flex items-center text-primary hover:text-primary-dark"
                variant="outline"
              >
                <PlusIcon className="h-5 w-5 mr-1" /> Add Category
              </Button>
            </PopoverTrigger>

            <PopoverContent sideOffset={30} className="w-full max-w-sm max-h-[70vh] left-0 mx-4 content-start overflow-y-auto p-4 bg-white">
              <div className="space-y-4">
                <h4 className="text-lg font-medium">New Category</h4>

                <div className="space-y-2">
                  <div className="grid grid-cols-3 items-center gap-4">
                    <Label htmlFor="cat-name">Name</Label>
                    <Input
                      id="cat-name"
                      value={newName}
                      onChange={(e) => setNewName(e.target.value)}
                      placeholder="E.g. work, fun"
                      className="col-span-2 h-8"
                    />
                  </div>

                  <div className="grid grid-cols-3 items-center gap-4">
                    <Label htmlFor="cat-color">Color</Label>
                    <Input
                      id="cat-color"
                      type="color"
                      value={newColor}
                      onChange={(e) => setNewColor(e.target.value)}
                      className="col-span-2 h-8 p-0"
                    />
                  </div>
                </div>

                {error && <p className="text-red-500 text-sm">{error}</p>}

                <div className="flex justify-end space-x-2">
                  <Button
                    variant="outline"
                    size="sm"
                    onClick={() => {
                      setNewName("");
                      setNewColor("#D69E2E");
                      setError(null);
                      setIsPopoverOpen(false);
                    }}
                  >
                    Cancel
                  </Button>
                  <Button size="sm" variant="outline" onClick={handleCreate} disabled={loading}>
                    {loading ? "Saving…" : "Save"}
                  </Button>
                </div>
              </div>
            </PopoverContent>
          </Popover>
        </div>
      </CardContent>
    </Card>
  );
}
