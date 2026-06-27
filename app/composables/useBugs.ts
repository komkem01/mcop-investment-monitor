import { ref, onMounted, onUnmounted } from "vue";

export interface BugItem {
  id: number; // Local index id
  sheetRowIndex: number; // The actual row index in Google Sheets (0-indexed)
  bugId: string;
  title: string;
  comments: string[];
  expectation: string;
  driveLink: string;
  figmaLink: string;
  status: string;
  createdDate: string;
  fixedDate: string;
  passedDate: string;
}

export function useBugs(sheetName: string) {
  const config = useRuntimeConfig();
  const API_BASE = config.public.apiBase;
  const bugsList = ref<BugItem[]>([]);
  const isLoading = ref(false);

  // ประกาศตัวแปรเก็บ Interval ไว้ที่นี่
  let pollInterval: ReturnType<typeof setInterval> | null = null;

  // Helper to parse formula: =HYPERLINK("url", "label") -> "url"
  const parseSheetLink = (val: any) => {
    if (!val) return "";
    const str = String(val);
    const match = str.match(/=HYPERLINK\("([^"]+)"/i) || str.match(/=HYPERLINK\('([^']+)'/i);
    return match ? match[1] : str;
  };

  // Helper to format URL: "url" -> `=HYPERLINK("url", "LINK")`
  const formatSheetLink = (url: string) => {
    if (!url || !url.trim()) return "";
    const cleanUrl = url.trim();
    if (cleanUrl.toLowerCase().startsWith("http")) {
      // Escape double quotes inside the URL just in case
      const escapedUrl = cleanUrl.replace(/"/g, '""');
      return `=HYPERLINK("${escapedUrl}", "LINK")`;
    }
    return cleanUrl;
  };

  // Fetch bugs from Google Sheets
  // เพิ่มพารามิเตอร์ showLoading เพื่อแยกว่าต้องขึ้นหน้าโหลดไหม
  const fetchBugs = async (showLoading = true) => {
    if (showLoading) {
      isLoading.value = true;
    }

    try {
      // Query up to 1000 records
      const response = await $fetch<any[][]>(API_BASE, {
        params: { range: `${sheetName}!A1:Z1000` },
      });

      if (!response || response.length < 3) {
        bugsList.value = [];
        return;
      }

      // Google sheets format:
      // Row 0: Metadata title
      // Row 1: Counts
      // Row 2: Headers
      // Row 3+: Actual data
      const dataRows = response.slice(3);

      bugsList.value = dataRows
        .map((row, index) => {
          if (!row || row.length === 0 || !row[0]) return null;

          // Row index in Google Sheets is index + 3 (since we sliced off first 3 rows)
          const sheetRowIndex = index + 3;

          // Parse comments list from newlines
          const rawComments = row[2] || "";
          const commentsList =
            typeof rawComments === "string"
              ? rawComments
                  .split("\n")
                  .map((c) => c.trim())
                  .filter((c) => c.length > 0)
              : [];

          return {
            id: index + 1,
            sheetRowIndex,
            bugId: String(row[0] || ""),
            title: String(row[1] || ""),
            comments: commentsList,
            expectation: String(row[3] || ""),
            driveLink: parseSheetLink(row[4]),
            figmaLink: parseSheetLink(row[5]),
            status: String(row[6] || "Defect"),
            createdDate: String(row[7] || ""),
            fixedDate: String(row[8] || ""),
            passedDate: String(row[9] || ""),
          } as BugItem;
        })
        .filter((b): b is BugItem => b !== null);
    } catch (error) {
      console.error("Failed to fetch bugs:", error);
    } finally {
      if (showLoading) {
        isLoading.value = false;
      }
    }
  };

  // Add a new bug
  const addBug = async (bug: Omit<BugItem, "id" | "sheetRowIndex">) => {
    const commentsStr = bug.comments.join("\n");
    const values = [
      [
        bug.bugId,
        bug.title,
        commentsStr,
        bug.expectation,
        formatSheetLink(bug.driveLink),
        formatSheetLink(bug.figmaLink),
        bug.status,
        bug.createdDate,
        bug.fixedDate,
        bug.passedDate,
      ],
    ];

    try {
      await $fetch(API_BASE, {
        method: "POST",
        body: {
          range: `${sheetName}!A4`,
          values,
        },
      });
      await fetchBugs(false); // ดึงข้อมูลใหม่แบบเงียบๆ ไม่ต้องแสดง Loading
      return true;
    } catch (error) {
      console.error("Failed to add bug:", error);
      return false;
    }
  };

  // Update an existing bug
  const updateBug = async (bug: BugItem) => {
    const commentsStr = bug.comments.join("\n");
    const values = [
      [
        bug.bugId,
        bug.title,
        commentsStr,
        bug.expectation,
        formatSheetLink(bug.driveLink),
        formatSheetLink(bug.figmaLink),
        bug.status,
        bug.createdDate,
        bug.fixedDate,
        bug.passedDate,
      ],
    ];

    try {
      // Row numbers in Google Sheets ranges are 1-indexed (e.g. A4:J4 for row index 3)
      const rangeStr = `${sheetName}!A${bug.sheetRowIndex + 1}:J${bug.sheetRowIndex + 1}`;
      await $fetch(API_BASE, {
        method: "PATCH",
        body: {
          range: rangeStr,
          values,
        },
      });
      await fetchBugs(false); // ดึงข้อมูลใหม่แบบเงียบๆ ไม่ต้องแสดง Loading
      return true;
    } catch (error) {
      console.error("Failed to update bug:", error);
      return false;
    }
  };

  // Delete a bug
  const deleteBug = async (sheetRowIndex: number) => {
    try {
      await $fetch(`${API_BASE}?sheetName=${encodeURIComponent(sheetName)}&rowIndex=${sheetRowIndex}`, {
        method: "DELETE",
      });
      await fetchBugs(false); // ดึงข้อมูลใหม่แบบเงียบๆ ไม่ต้องแสดง Loading
      return true;
    } catch (error) {
      console.error("Failed to delete bug:", error);
      return false;
    }
  };

  // ตั้งค่าการดึงข้อมูลอัตโนมัติเมื่อ Component ถูกเรียกใช้งาน
  onMounted(() => {
    // 120 วินาที = 120,000 มิลลิวินาที
    pollInterval = setInterval(() => {
      fetchBugs(false); // ดึงข้อมูลพื้นหลัง ไม่โชว์ Loading
    }, 120000);
  });

  // เคลียร์รอบการทำงานเมื่อ Component ถูกทำลาย (ปิดหน้า)
  onUnmounted(() => {
    if (pollInterval) {
      clearInterval(pollInterval);
    }
  });

  return {
    bugsList,
    isLoading,
    fetchBugs,
    addBug,
    updateBug,
    deleteBug,
  };
}
