import { UserPlus } from 'lucide-react'
import { AuthenticatedLayout } from '@/layouts/authenticated'
import { Button } from '@/shared/components/ui/button'
import { useUsers } from './hooks/use-users'
import { useInvites } from './hooks/use-invites'
import { UsersTable } from './ui/users-table'
import { PendingInvitesTable } from './ui/pending-invites-table'
import { InviteUserDialog } from './ui/invite-user-dialog'
import { InviteLinkDialog } from './ui/invite-link-dialog'
import { ChangeRoleDialog } from './ui/change-role-dialog'
import { RemoveUserDialog } from './ui/remove-user-dialog'
import { CancelInviteDialog } from './ui/cancel-invite-dialog'

export function UsersPage() {
  const {
    users,
    isLoading: isLoadingUsers,
    changeRoleDialogOpen,
    removeUserDialogOpen,
    selectedUser,
    isUpdatingRole,
    isRemovingUser,
    handleChangeRole,
    confirmChangeRole,
    handleRemoveUser,
    confirmRemoveUser,
    closeDialogs: closeUserDialogs,
  } = useUsers()

  const {
    invites,
    isLoading: isLoadingInvites,
    inviteDialogOpen,
    setInviteDialogOpen,
    inviteLinkDialogOpen,
    inviteLink,
    cancelInviteDialogOpen,
    selectedInvite,
    isCreatingInvite,
    isCancellingInvite,
    handleCreateInvite,
    handleCancelInvite,
    confirmCancelInvite,
    closeDialogs: closeInviteDialogs,
  } = useInvites()

  return (
    <AuthenticatedLayout
      title="Users"
      actions={
        <Button onClick={() => setInviteDialogOpen(true)}>
          <UserPlus className="h-4 w-4" />
          Invite User
        </Button>
      }
    >
      <div className="space-y-8">
        <PendingInvitesTable
          invites={invites}
          isLoading={isLoadingInvites}
          onCancel={handleCancelInvite}
        />

        <UsersTable
          users={users}
          isLoading={isLoadingUsers}
          onChangeRole={handleChangeRole}
          onRemove={handleRemoveUser}
        />
      </div>

      <InviteUserDialog
        open={inviteDialogOpen}
        isLoading={isCreatingInvite}
        onInvite={handleCreateInvite}
        onClose={closeInviteDialogs}
      />

      <InviteLinkDialog
        open={inviteLinkDialogOpen}
        inviteLink={inviteLink}
        onClose={closeInviteDialogs}
      />

      <CancelInviteDialog
        open={cancelInviteDialogOpen}
        invite={selectedInvite}
        isLoading={isCancellingInvite}
        onConfirm={confirmCancelInvite}
        onClose={closeInviteDialogs}
      />

      <ChangeRoleDialog
        open={changeRoleDialogOpen}
        user={selectedUser}
        isLoading={isUpdatingRole}
        onConfirm={confirmChangeRole}
        onClose={closeUserDialogs}
      />

      <RemoveUserDialog
        open={removeUserDialogOpen}
        user={selectedUser}
        isLoading={isRemovingUser}
        onConfirm={confirmRemoveUser}
        onClose={closeUserDialogs}
      />
    </AuthenticatedLayout>
  )
}
