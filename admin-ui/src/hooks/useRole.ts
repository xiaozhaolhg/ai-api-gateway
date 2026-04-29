import { useAuth } from '../contexts/AuthContext';

type Role = 'admin' | 'user' | 'viewer';

export function useRole(): Role {
  const { user } = useAuth();
  return (user?.role || 'viewer') as Role;
}
