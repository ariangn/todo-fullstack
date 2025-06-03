import { useEffect, useState } from "react";
import {
  Select,
  SelectTrigger,
  SelectValue,
  SelectContent,
  SelectItem,
} from "@/components/ui/select";
import {
  DropdownMenu,
  DropdownMenuTrigger,
  DropdownMenuContent,
  DropdownMenuGroup,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
} from "@/components/ui/dropdown-menu";
import { Button } from "@/components/ui/button";
import { ChevronDownIcon, FilterIcon } from "lucide-react";
import { fetchCategories } from "../services/categoryService";
import { fetchTags } from "../services/tagService";

// Define a union type for the allowed sort keys:
type SortKey = "dueDate" | "createdAt" | "updatedAt";

export default function SortFilterBar({
  onSortChange,
  onFilterChange,
}: {
  onSortChange: (sortBy: SortKey) => void;
  onFilterChange: (selectedCategoryIds: string[], selectedTagIds: string[]) => void;
}) {
  const [categories, setCategories] = useState<{ id: string; name: string }[]>([]);
  const [tags, setTags] = useState<{ id: string; name: string }[]>([]);
  const [selectedCats, setSelectedCats] = useState<string[]>([]);
  const [selectedTags, setSelectedTags] = useState<string[]>([]);

  useEffect(() => {
    fetchCategories().then((cats) => setCategories(cats));
    fetchTags().then((t) => setTags(t));
  }, []);

  const toggleCategory = (id: string) => {
    const updated = selectedCats.includes(id)
      ? selectedCats.filter((x) => x !== id)
      : [...selectedCats, id];
    setSelectedCats(updated);
    onFilterChange(updated, selectedTags);
  };

  const toggleTag = (id: string) => {
    const updated = selectedTags.includes(id)
      ? selectedTags.filter((x) => x !== id)
      : [...selectedTags, id];
    setSelectedTags(updated);
    onFilterChange(selectedCats, updated);
  };

  return (
    <div className="flex items-center space-x-4 mb-4">
      <Select
        onValueChange={(v: string) =>
          onSortChange(v as SortKey)
        }
      >
        <SelectTrigger className="w-44">
          <SelectValue placeholder="Sort by" />
          <ChevronDownIcon className="ml-auto h-4 w-4" />
        </SelectTrigger>
        <SelectContent>
          <SelectItem value="dueDate">Due Date</SelectItem>
          <SelectItem value="createdAt">Date Created</SelectItem>
          <SelectItem value="updatedAt">Date Updated</SelectItem>
        </SelectContent>
      </Select>

      <DropdownMenu>
        <DropdownMenuTrigger asChild>
          <Button variant="outline" className="flex items-center space-x-2">
            <FilterIcon className="h-4 w-4" />
            <span>Filter</span>
            <ChevronDownIcon className="h-4 w-4 ml-1" />
          </Button>
        </DropdownMenuTrigger>
        <DropdownMenuContent className="w-64">
          <DropdownMenuLabel className="font-medium">Categories</DropdownMenuLabel>
          <DropdownMenuGroup>
            {categories.map((cat) => (
              <DropdownMenuItem
                key={cat.id}
                className="flex items-center space-x-2"
                onClick={() => toggleCategory(cat.id)}
              >
                <input
                  type="checkbox"
                  checked={selectedCats.includes(cat.id)}
                  readOnly
                  className="form-checkbox h-4 w-4 text-primary"
                />
                <span>{cat.name}</span>
              </DropdownMenuItem>
            ))}
          </DropdownMenuGroup>
          <DropdownMenuSeparator />
          <DropdownMenuLabel className="font-medium">Tags</DropdownMenuLabel>
          <DropdownMenuGroup>
            {tags.map((tg) => (
              <DropdownMenuItem
                key={tg.id}
                className="flex items-center space-x-2"
                onClick={() => toggleTag(tg.id)}
              >
                <input
                  type="checkbox"
                  checked={selectedTags.includes(tg.id)}
                  readOnly
                  className="form-checkbox h-4 w-4 text-primary"
                />
                <span>{tg.name}</span>
              </DropdownMenuItem>
            ))}
          </DropdownMenuGroup>
        </DropdownMenuContent>
      </DropdownMenu>
    </div>
  );
}
