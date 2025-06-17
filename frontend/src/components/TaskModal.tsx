// src/components/TaskModal.tsx
import { useEffect, useState } from "react";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogFooter,
  DialogDescription,
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
  /** FULL list of all tags in the system */
  allTags: { id: string; name: string }[];
  /** createTag should upsert by name and return { id, name } */
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
  allTags,
  createTag,
  onSave,
  onDelete,
  onClose,
}: TaskModalProps) {
  const [title, setTitle] = useState(todo?.title || "");
  const [body, setBody] = useState(todo?.body || "");
  const [dueDate, setDueDate] = useState<Date | undefined>(
    todo?.dueDate ? new Date(todo.dueDate) : undefined
  );
  const [categoryId, setCategoryId] = useState<string | undefined>(
    todo?.categoryId || undefined
  );
  const [tagInput, setTagInput] = useState("");
  /** Now an array of full tag objects */
  const [selectedTags, setSelectedTags] = useState<
    { id: string; name: string }[]
  >([]);
  const [tagError, setTagError] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  // initialize on open / edit
  useEffect(() => {
    if (mode === "edit" && todo) {
      setTitle(todo.title);
      setBody(todo.body || "");
      setDueDate(todo.dueDate ? new Date(todo.dueDate) : undefined);
      setCategoryId(todo.categoryId || undefined);
      // match IDs → full objects
      setSelectedTags(
        allTags.filter((t) => todo.tagIds.includes(t.id))
      );
    } else {
      // create mode: clear
      setTitle("");
      setBody("");
      setDueDate(undefined);
      setCategoryId(undefined);
      setSelectedTags([]);
    }
  }, [mode, todo, allTags]);

  const handleAddTag = async () => {
    const raw = tagInput.trim();
    // 1) Validate non-empty
    if (!raw) {
      setTagError("Tag cannot be empty");
      return;
    }
    // 2) Validate no duplicates
    if (selectedTags.some((t) => t.name.toLowerCase() === raw.toLowerCase())) {
      setTagError(`"${raw}" is already added`);
      return;
    }

    setTagError(null);
    setLoading(true);
    try {
      const tag = await createTag(raw);
      setSelectedTags((prev) => [...prev, tag]);
    } catch (e) {
      console.error("Failed to create/tag:", e);
      setTagError("Failed to add tag");
    } finally {
      setLoading(false);
      setTagInput("");
    }
  };

  const handleRemoveTag = (id: string) => {
    setSelectedTags((prev) => prev.filter((t) => t.id !== id));
  };

  const handleSave = async () => {
    if (!title.trim()) {
      setError("Title is required");
      return;
    }
    setError(null);
    setLoading(true);

    try {
      const tagIds = selectedTags.map((t) => t.id);
      await onSave({
        title: title.trim(),
        body: body.trim() || undefined,
        dueDate: dueDate ? dueDate.toISOString() : undefined,
        status: mode === "create" ? status! : todo!.status,
        categoryId,
        tagIds,
      });
      onClose();
    } catch (err: any) {
      setError(err.message || String(err));
    } finally {
      setLoading(false);
    }
  };

  const handleDelete = async () => {
    if (!onDelete) return;
    setLoading(true);
    try {
      await onDelete();
      onClose();
    } catch (err: any) {
      setError(err.message || String(err));
    } finally {
      setLoading(false);
    }
  };

  return (
    <Dialog open onOpenChange={onClose}>
      <DialogContent className="sm:max-w-lg bg-white">
        <DialogHeader>
          <DialogTitle>{mode === "create" ? "New Task" : "Edit Task"}</DialogTitle>
          <DialogDescription>Fill out the details below.</DialogDescription>
        </DialogHeader>

        <div className="space-y-4 mt-2">
          {/* Title, Body, Due Date, Category… */}
          <div className="grid grid-cols-3 items-center gap-4">
            <Label htmlFor="task-title">Title</Label>
            <Input
              id="task-title"
              value={title}
              onChange={(e) => setTitle(e.target.value)}
              className="col-span-2 h-8"
            />
          </div>
          {/* Notes */}
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
          {/* Due Date */}
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
          {/* Category */}
          <div className="grid grid-cols-3 items-center gap-4">
            <Label htmlFor="task-category">Category</Label>
            <Select
              onValueChange={(v: string) => setCategoryId(v === "none" ? undefined : v)}
              value={categoryId}
            >
              <SelectTrigger className="w-[180px]">
                <SelectValue placeholder="Select category" />
              </SelectTrigger>
              <SelectContent className="">
                <SelectItem value="none">None</SelectItem>
                {categories.map((cat) => (
                  <SelectItem key={cat.id} value={cat.id}>
                    {cat.name}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>

          {/* Tags */}
          <div className="grid grid-cols-3 items-start gap-4">
            <Label>Tags</Label>
            <div className="col-span-2 space-y-2">
              <div className="flex space-x-2">
                <Input
                  value={tagInput}
                  onChange={(e) => {
                    setTagInput(e.target.value);
                    if (tagError) setTagError(null);
                  }}
                  placeholder="Add a tag"
                  className="flex-1 h-8"
                />
                <Button size="sm" onClick={handleAddTag} disabled={loading}>
                  +
                </Button>
              </div>
              {/* show tag-specific error */}
              {tagError && (
                <p className="text-red-500 text-xs">{tagError}</p>
              )}
              <div className="flex flex-wrap gap-2">
                {selectedTags.map((t) => (
                  <Badge
                    key={t.id}
                    className="flex items-center space-x-1 cursor-pointer"
                  >
                    <span>{t.name}</span>
                    <button
                      onClick={() => handleRemoveTag(t.id)}
                      className="p-1"
                    >
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
              variant="outline"
              onClick={handleDelete}
              disabled={loading}
            >
              Delete
            </Button>
          )}
          <Button variant="outline" onClick={handleSave} disabled={loading}>
            {loading ? "Saving…" : "Save"}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
