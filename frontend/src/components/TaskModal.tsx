// src/components/TaskModal.tsx
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
import { Textarea } from "@/components/ui/textarea";
import {
  Select,
  SelectTrigger,
  SelectValue,
  SelectContent,
  SelectItem,
} from "@/components/ui/select";
import { Popover, PopoverTrigger, PopoverContent } from "@/components/ui/popover";
import { Calendar } from "@/components/ui/calendar";
import { Badge } from "@/components/ui/badge";
import { XIcon } from "lucide-react";
import type { Todo } from "../services/todoService";
import type { Category } from "../types";

interface TaskModalProps {
  mode: "create" | "edit";
  status?: Todo["status"];
  todo?: Todo;
  categories: Category[];
  tagIds: string[];  
  createTag: (name: string) => Promise<{ id: string; name: string }>;
  onSave: (data: {
    title: string;
    body?: string;
    dueDate?: string;
    status: Todo["status"];
    categoryId?: string;
    tagIds: string[];
  }) => Promise<void>;
  onDelete?: () => Promise<void>;
  onClose: () => void;
}

export default function TaskModal({
  mode,
  status,
  todo,
  categories,
  createTag,
  onSave,
  onDelete,
  onClose,
}: TaskModalProps) {
  const [title, setTitle] = useState<string>(todo?.title || "");
  const [body, setBody] = useState<string>(todo?.body || "");
  const [dueDate, setDueDate] = useState<Date | undefined>(
    todo?.dueDate ? new Date(todo.dueDate) : undefined
  );
  const [categoryId, setCategoryId] = useState<string | undefined>(
    todo?.categoryId ?? undefined
  );
  const [tagInput, setTagInput] = useState<string>("");
  const [selectedTags, setSelectedTags] = useState<string[]>(todo?.tagIds || []);
  const [loading, setLoading] = useState<boolean>(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const badCats = categories.filter((cat) => !cat.id || cat.id.trim() === "");
    if (badCats.length > 0) {
      console.warn("Invalid categories detected:", badCats);
    }
    if (mode === "edit" && todo) {
      setTitle(todo.title);
      setBody(todo.body || "");
      setDueDate(todo.dueDate ? new Date(todo.dueDate) : undefined);
      setCategoryId(todo.categoryId || "");
      setSelectedTags(todo.tagIds);
    }
  }, [mode, todo, categories]);

  const handleAddTag = () => {
    const trimmed = tagInput.trim();
    if (!trimmed) return;
    if (!selectedTags.includes(trimmed)) {
      setSelectedTags([...selectedTags, trimmed]);
    }
    setTagInput("");
  };

  const handleRemoveTag = (tagName: string) => {
    setSelectedTags(selectedTags.filter((t) => t !== tagName));
  };

  const handleSave = async () => {
    if (!title.trim()) {
      setError("Title is required");
      return;
    }

    setError(null);
    setLoading(true);

    try {
      const tagIds: string[] = [];
      console.log("selectedTags:", selectedTags);
      for (const tagName of selectedTags) {
        const tag = await createTag(tagName); // create or get existing
        tagIds.push(tag.id);
      }

      await onSave({
        title: title.trim(),
        body: body.trim() || undefined,
        dueDate: dueDate ? dueDate.toISOString() : undefined,
        status: mode === "create" ? status! : todo!.status,
        categoryId: categoryId !== "none" ? categoryId : undefined,
        tagIds: tagIds,
      });

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
      <DialogContent className="sm:max-w-lg">
        <DialogHeader>
          <DialogTitle>
            {mode === "create" ? "New Task" : "Edit Task"}
          </DialogTitle>
        </DialogHeader>

        <div className="space-y-4 mt-2">
          <div className="grid grid-cols-3 items-center gap-4">
            <Label htmlFor="task-title">Title</Label>
            <Input
              id="task-title"
              value={title}
              onChange={(e) => setTitle(e.target.value)}
              className="col-span-2 h-8"
            />
          </div>

          <div className="grid grid-cols-3 items-start gap-4">
            <Label htmlFor="task-body">Notes</Label>
            <Textarea
              id="task-body"
              value={body}
              onChange={(e) => setBody(e.target.value)}
              className="col-span-2"
              rows={3}
            />
          </div>

          <div className="grid grid-cols-3 items-center gap-4">
            <Label>Due Date</Label>
            <Popover>
              <PopoverTrigger asChild>
                <Button
                  variant="outline"
                  className="w-[200px] justify-start text-left font-normal"
                >
                  {dueDate
                    ? dueDate.toLocaleDateString(undefined, {
                        year: "numeric",
                        month: "short",
                        day: "numeric",
                      })
                    : "Select date"}
                </Button>
              </PopoverTrigger>
              <PopoverContent className="w-auto p-0">
                <Calendar
                  mode="single"
                  selected={dueDate}
                  onSelect={setDueDate}
                  initialFocus
                />
              </PopoverContent>
            </Popover>
          </div>

          <div className="grid grid-cols-3 items-center gap-4">
            <Label htmlFor="task-category">Category</Label>
            <Select
              onValueChange={(v: string) => setCategoryId(v)}
              value={categoryId ?? undefined}
            >
              <SelectTrigger className="w-[180px]">
                <SelectValue placeholder="Select category" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="none">None</SelectItem>
                {categories.map((cat) => (
                  <SelectItem key={cat.id} value={cat.id}>
                    {cat.name}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>

          <div className="grid grid-cols-3 items-center gap-4">
            <Label htmlFor="task-tags">Tags</Label>
            <div className="col-span-2 space-y-2">
              <div className="flex space-x-2">
                <Input
                  id="task-tags"
                  value={tagInput}
                  onChange={(e) => setTagInput(e.target.value)}
                  placeholder="Add tag"
                  className="flex-1 h-8"
                />
                <Button size="sm" onClick={handleAddTag}>
                  +
                </Button>
              </div>
              <div className="flex flex-wrap space-x-1">
                {selectedTags.map((tagName) => (
                  <Badge key={tagName} className="flex items-center space-x-1">
                    <span>{tagName}</span>
                    <button onClick={() => handleRemoveTag(tagName)}>
                      <XIcon className="h-4 w-4" />
                    </button>
                  </Badge>
                ))}
              </div>
            </div>
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
