import React from 'react';
import { Route } from 'react-router-dom';

import { Registration } from './auth/registration';

export const CustomRoutes = [
  <Route key="register" path="/register" component={Registration} noLayout />,
];
