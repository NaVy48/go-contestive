import React from 'react';
import { Login } from 'react-admin';

import { RegistrationForm } from './registrationForm';

export const Registration = () => {
  return (
    <Login>
      <RegistrationForm />
    </Login>
  );
};
