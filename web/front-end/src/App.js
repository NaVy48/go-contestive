import './App.css';

import { LoginPage } from 'components/auth/login';
import { ContestsComponents, UserContestsComponents } from 'components/contests';
import { CustomRoutes } from 'components/customRoutes';
import { ProblemsComponents, UserProblemsComponents } from 'components/problems';
import { SubmissionsComponents } from 'components/submissions';
import { UsersComponents } from 'components/users';
import React from 'react';
import { Admin, Resource } from 'react-admin';
import { Redirect } from 'react-router'

import { authProvider } from './providers/authProvider';
import { dataProvider } from './providers/dataProvider';

export const App = () => {
  return (
    <Admin
      dataProvider={dataProvider}
      authProvider={authProvider}
      customRoutes={CustomRoutes}
      loginPage={LoginPage}
    >
      {permissions =>
        permissions === 'admin'
          ? [
            <Resource key="users" name="users" {...UsersComponents} />,
            <Resource key="problems" name="problems" {...ProblemsComponents} />,
            <Resource key="contests" name="contests" {...ContestsComponents} />,
            <Resource key="submissions" name="submissions" {...SubmissionsComponents} />,
            <Resource key="standings" name="standings" />,
          ]
          : permissions === 'user'
            ? [
              <Resource key="contests" name="contests" {...UserContestsComponents} />,
              <Resource key="problems" name="problems" {...UserProblemsComponents} />,
              <Resource key="submissions" name="submissions" {...SubmissionsComponents} />,
              <Resource key="standings" name="standings" />,
            ]
            : [
              <Redirect to="/login" />
            ]
      }
    </Admin>
  );
};
