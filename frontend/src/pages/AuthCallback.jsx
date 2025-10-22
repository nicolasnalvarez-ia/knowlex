import { useEffect } from 'react';
import { useNavigate, useSearchParams } from 'react-router-dom';
import useAuthStore from '../lib/authStore';

const AuthCallback = () => {
  const navigate = useNavigate();
  const [searchParams] = useSearchParams();
  const { setAuth, fetchUser } = useAuthStore();

  useEffect(() => {
    const token = searchParams.get('token');
    
    if (token) {
      localStorage.setItem('token', token);
      
      fetchUser()
        .then(() => {
          navigate('/dashboard', { replace: true });
        })
        .catch(() => {
          navigate('/', { replace: true });
        });
    } else {
      navigate('/', { replace: true });
    }
  }, [searchParams, navigate, setAuth, fetchUser]);

  return (
    <div className="min-h-screen bg-gray-100 dark:bg-gray-900 flex items-center justify-center">
      <div className="text-center">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-500 mx-auto mb-4"></div>
        <p className="text-gray-600 dark:text-gray-300">Logging you in...</p>
      </div>
    </div>
  );
};

export default AuthCallback;
