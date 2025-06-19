// src/components/SortFilterBar.tsx
import { useEffect, useState, useMemo } from "react";
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
  const [selectedTagNames, setSelectedTagNames] = useState<string[]>([]);

  // Fetch once
  useEffect(() => {
    fetchCategories().then((cats) => setCategories(cats || []));
    fetchTags().then((t) => setTags(t || []));
  }, []);

  // Build a map: tagName → [tagId, ...]
  const nameToIds = useMemo(() => {
    const m = new Map<string, string[]>();
    tags.forEach(({ id, name }) => {
      const arr = m.get(name) ?? [];
      arr.push(id);
      m.set(name, arr);
    });
    return m;
  }, [tags]);

  // Unique list of names for the UI
  const uniqueNames = useMemo(() => Array.from(nameToIds.keys()), [nameToIds]);

  // toggle a category, then call parent
  const toggleCategory = (id: string) => {
    const updated = selectedCats.includes(id)
      ? selectedCats.filter((x) => x !== id)
      : [...selectedCats, id];
    setSelectedCats(updated);
    onFilterChange(updated, flattenTagIds(selectedTagNames));
  };

  // toggle a tag name, then call parent
  const toggleTagName = (name: string) => {
    const updatedNames = selectedTagNames.includes(name)
      ? selectedTagNames.filter((x) => x !== name)
      : [...selectedTagNames, name];
    setSelectedTagNames(updatedNames);
    onFilterChange(selectedCats, flattenTagIds(updatedNames));
  };

  function flattenTagIds(names: string[]) {
    return names.flatMap((n) => nameToIds.get(n) ?? []);
  }

  return (
    <div className="flex items-center space-x-4 mb-4">
      {/* Sort dropdown */}
      <Select onValueChange={(v) => onSortChange(v as SortKey)}>
        <SelectTrigger className="w-44 hover:bg-gray-300">
          <SelectValue placeholder="Sort by" />
        </SelectTrigger>
        <SelectContent className="bg-white">
          <SelectItem value="dueDate">Due Date</SelectItem>
          <SelectItem value="createdAt">Date Created</SelectItem>
          <SelectItem value="updatedAt">Date Updated</SelectItem>
        </SelectContent>
      </Select>

      {/* Filter menu */}
      <DropdownMenu>
        <DropdownMenuTrigger asChild>
          <Button variant="outline" className="flex items-center space-x-2">
            <FilterIcon className="h-4 w-4" />
            <span>Filter</span>
            <ChevronDownIcon className="h-4 w-4 ml-1" />
          </Button>
        </DropdownMenuTrigger>
        <DropdownMenuContent className="w-64 bg-white">

          {/* Categories */}
          <DropdownMenuLabel className="font-medium">Categories</DropdownMenuLabel>
          <DropdownMenuGroup>
            {categories.length > 0 ? (
              categories.map((cat) => (
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
              ))
            ) : (
              <DropdownMenuItem disabled className="text-sm italic text-gray-500">
                No categories created yet
              </DropdownMenuItem>
            )}
          </DropdownMenuGroup>

          <DropdownMenuSeparator />

          {/* Tags (deduped by name) */}
          <DropdownMenuLabel className="font-medium">Tags</DropdownMenuLabel>
          <DropdownMenuGroup>
            {uniqueNames.length > 0 ? (
              uniqueNames.map((name) => (
                <DropdownMenuItem
                  key={name}
                  className="flex items-center space-x-2"
                  onClick={() => toggleTagName(name)}
                >
                  <input
                    type="checkbox"
                    checked={selectedTagNames.includes(name)}
                    readOnly
                    className="form-checkbox h-4 w-4 text-primary"
                  />
                  <span>{name}</span>
                </DropdownMenuItem>
              ))
            ) : (
              <DropdownMenuItem disabled className="text-sm italic text-gray-500">
                No tags created yet
              </DropdownMenuItem>
            )}
          </DropdownMenuGroup>

        </DropdownMenuContent>
      </DropdownMenu>
    </div>
  );

}
