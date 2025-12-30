import React from 'react';
import { useNavigate } from 'react-router';

export default function IndexPage() {
  const navigate = useNavigate();

  React.useEffect(() => {
    const token = localStorage.getItem('token');
    if (token) {
      navigate('/dashboard/users');
    } else {
      navigate('/login');
    }
  }, [navigate]);

  return (
    <div className="min-h-screen bg-gradient-to-br from-blue-600 to-blue-800 flex items-center justify-center">
      <div className="text-center">
        <div className="animate-spin rounded-full h-16 w-16 border-b-2 border-white mx-auto"></div>
        <p className="mt-4 text-white text-xl">Loading...</p>
      </div>
    </div>
  );
}
