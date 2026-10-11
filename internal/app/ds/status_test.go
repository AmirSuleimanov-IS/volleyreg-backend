package ds

import "testing"

func TestCanChangeStatus(t *testing.T) {
	cases := []struct {
		from, to    string
		byModerator bool
		want        bool
	}{
		// представитель команды удаляет и формирует черновик
		{StatusDraft, StatusFormed, false, true},
		{StatusDraft, StatusDeleted, false, true},
		// организатор допускает или отклоняет сформированную заявку
		{StatusFormed, StatusCompleted, true, true},
		{StatusFormed, StatusRejected, true, true},
		// всё остальное запрещено
		{StatusDraft, StatusCompleted, true, false},    // завершить черновик нельзя
		{StatusFormed, StatusCompleted, false, false},  // представитель не может допустить заявку
		{StatusDraft, StatusFormed, true, false},       // организатор не формирует чужой черновик
		{StatusFormed, StatusDeleted, false, false},    // сформированную заявку не удалить
		{StatusCompleted, StatusRejected, true, false}, // завершённую не изменить
		{StatusDeleted, StatusFormed, false, false},
	}
	for _, c := range cases {
		if got := CanChangeStatus(c.from, c.to, c.byModerator); got != c.want {
			t.Errorf("CanChangeStatus(%s -> %s, moderator=%v) = %v, want %v", c.from, c.to, c.byModerator, got, c.want)
		}
	}
}
