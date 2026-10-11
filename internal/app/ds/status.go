package ds

// CanChangeStatus — можно ли перевести заявку из статуса from в статус to.
// Представитель команды удаляет и формирует черновик,
// организатор соревнования допускает или отклоняет сформированную заявку.
func CanChangeStatus(from, to string, byModerator bool) bool {
	if byModerator {
		return from == StatusFormed && (to == StatusCompleted || to == StatusRejected)
	}
	return from == StatusDraft && (to == StatusFormed || to == StatusDeleted)
}
