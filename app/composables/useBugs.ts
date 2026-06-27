import { ref } from 'vue'

const API_BASE = 'http://localhost:8080/api/data'

export interface BugItem {
  id: number          // Local index id
  sheetRowIndex: number // The actual row index in Google Sheets (0-indexed)
  bugId: string
  title: string
  comments: string[]
  expectation: string
  driveLink: string
  figmaLink: string
  status: string
  createdDate: string
  fixedDate: string
  passedDate: string
}

export function useBugs(sheetName: string) {
  const bugsList = ref<BugItem[]>([])
  const isLoading = ref(false)

  // Fetch bugs from Google Sheets
  const fetchBugs = async () => {
    isLoading.value = true
    try {
      // Query up to 1000 records
      const response = await $fetch<any[][]>(API_BASE, {
        params: { range: `${sheetName}!A1:Z1000` }
      })

      if (!response || response.length < 3) {
        bugsList.value = []
        return
      }

      // Google sheets format:
      // Row 0: Metadata title
      // Row 1: Counts
      // Row 2: Headers
      // Row 3+: Actual data
      const dataRows = response.slice(3)

      bugsList.value = dataRows.map((row, index) => {
        if (!row || row.length === 0 || !row[0]) return null

        // Row index in Google Sheets is index + 3 (since we sliced off first 3 rows)
        const sheetRowIndex = index + 3

        // Parse comments list from newlines
        const rawComments = row[2] || ''
        const commentsList = typeof rawComments === 'string'
          ? rawComments.split('\n').map(c => c.trim()).filter(c => c.length > 0)
          : []

        return {
          id: index + 1,
          sheetRowIndex,
          bugId: String(row[0] || ''),
          title: String(row[1] || ''),
          comments: commentsList,
          expectation: String(row[3] || ''),
          driveLink: String(row[4] || ''),
          figmaLink: String(row[5] || ''),
          status: String(row[6] || 'Defect'),
          createdDate: String(row[7] || ''),
          fixedDate: String(row[8] || ''),
          passedDate: String(row[9] || '')
        } as BugItem
      }).filter((b): b is BugItem => b !== null)

    } catch (error) {
      console.error('Failed to fetch bugs:', error)
    } finally {
      isLoading.value = false
    }
  }

  // Add a new bug
  const addBug = async (bug: Omit<BugItem, 'id' | 'sheetRowIndex'>) => {
    const commentsStr = bug.comments.join('\n')
    const values = [
      [
        bug.bugId,
        bug.title,
        commentsStr,
        bug.expectation,
        bug.driveLink,
        bug.figmaLink,
        bug.status,
        bug.createdDate,
        bug.fixedDate,
        bug.passedDate
      ]
    ]

    try {
      await $fetch(API_BASE, {
        method: 'POST',
        body: {
          range: `${sheetName}!A4`,
          values
        }
      })
      await fetchBugs() // Refresh list
      return true
    } catch (error) {
      console.error('Failed to add bug:', error)
      return false
    }
  }

  // Update an existing bug
  const updateBug = async (bug: BugItem) => {
    const commentsStr = bug.comments.join('\n')
    const values = [
      [
        bug.bugId,
        bug.title,
        commentsStr,
        bug.expectation,
        bug.driveLink,
        bug.figmaLink,
        bug.status,
        bug.createdDate,
        bug.fixedDate,
        bug.passedDate
      ]
    ]

    try {
      // Row numbers in Google Sheets ranges are 1-indexed (e.g. A4:J4 for row index 3)
      const rangeStr = `${sheetName}!A${bug.sheetRowIndex + 1}:J${bug.sheetRowIndex + 1}`
      await $fetch(API_BASE, {
        method: 'PUT',
        body: {
          range: rangeStr,
          values
        }
      })
      await fetchBugs() // Refresh list
      return true
    } catch (error) {
      console.error('Failed to update bug:', error)
      return false
    }
  }

  // Delete a bug
  const deleteBug = async (sheetRowIndex: number) => {
    try {
      await $fetch(`${API_BASE}?sheetName=${encodeURIComponent(sheetName)}`, {
        method: 'DELETE',
        body: {
          rowIndex: sheetRowIndex
        }
      })
      await fetchBugs() // Refresh list
      return true
    } catch (error) {
      console.error('Failed to delete bug:', error)
      return false
    }
  }

  return {
    bugsList,
    isLoading,
    fetchBugs,
    addBug,
    updateBug,
    deleteBug
  }
}
