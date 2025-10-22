import { Navigate } from 'react-router-dom';
import useAuthStore from '../lib/authStore';

export const useAuth = () => {
  const { isAuthenticated, user } = useAuthStore();
  return { isAuthenticated, user };
};

export const ProtectedRoute = ({ children }) => {
  const { isAuthenticated } = useAuth();
  
  if (!isAuthenticated) {
    return <Navigate to="/" replace />;
  }
  
  return children;
};
