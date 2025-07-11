// src/components/ui/calendar-24.tsx
"use client"

import * as React from "react"
import { useState, useEffect } from "react"
import { ChevronDownIcon } from "lucide-react"

import { Button } from "@/components/ui/button"
import { Calendar } from "@/components/ui/calendar"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "@/components/ui/popover"

interface Calendar24Props {
  date?: Date
  onDateChange: (date?: Date) => void
  time?: Date
  onTimeChange: (time: string) => void
}

export function Calendar24({
  date,
  onDateChange,
  time,
  onTimeChange,
}: Calendar24Props) {
  const [open, setOpen] = React.useState(false)
  const [timeStr, setTimeStr] = useState("");

  // Sync local state whenever the prop changes
  useEffect(() => {
    if (time) {
      // format as HH:MM:SS
      setTimeStr(time.toTimeString().slice(0, 8));
    } else {
      setTimeStr("");
    }
  }, [time]);

  return (
    <div className="flex gap-4">
      <div className="flex flex-col gap-3">
        <Label htmlFor="date-picker" className="px-1">
          Date
        </Label>
        <Popover open={open} onOpenChange={setOpen}>
          <PopoverTrigger asChild>
            <Button
              variant="outline"
              id="date-picker"
              className="w-32 justify-between font-normal"
            >
              {date ? date.toLocaleDateString() : "Select date"}
              <ChevronDownIcon />
            </Button>
          </PopoverTrigger>
          <PopoverContent className="w-auto overflow-hidden p-0 bg-white" align="start">
            <Calendar
              mode="single"
              selected={date}
              captionLayout="dropdown"
              onSelect={(d) => {
                onDateChange(d)
                setOpen(false)
              }}
            />
          </PopoverContent>
        </Popover>
      </div>
      <div className="flex flex-col gap-3">
        <Label htmlFor="time-picker" className="px-1">
          Time
        </Label>
        <Input
          type="time"
          id="time-picker"
          step="1"
          value={timeStr}
          onChange={(e) => {
            setTimeStr(e.target.value);       // update local only
          }}
          onBlur={() => {
            // only commit when the string is a full HH:MM:SS
            const parts = timeStr.split(":").map(Number);
            if (parts.length === 3 && parts.every((n) => !isNaN(n))) {
              onTimeChange(timeStr);
            }
          }}
          className="bg-background appearance-none [&::-webkit-calendar-picker-indicator]:hidden [&::-webkit-calendar-picker-indicator]:appearance-none"
        />
      </div>
    </div>
  )
}
