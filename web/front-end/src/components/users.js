import PeopleIcon from '@material-ui/icons/People';
import React from 'react';
import {
  BooleanField,
  BooleanInput,
  Create,
  Datagrid,
  DateField,
  List,
  PasswordInput,
  SimpleForm,
  TextField,
  TextInput,
  EditButton,
  Edit,
  Show,
  SimpleShowLayout,
} from 'react-admin';



const UserTitle = ({ label = "User", record }) => {
  return <span>{label} {record ? `"${record.username}"` : ''}</span>;
};


export const UserList = props => (
  <List {...props} exporter={false} bulkActionButtons={false}>
    <Datagrid rowClick="show">
      <TextField source="id" />
      <TextField source="username" />
      <TextField source="firstName" label="First Name" />
      <TextField source="lastName" label="Last Name" />
      <DateField source="createdAt" label="Created At" />
      <BooleanField source="admin" />
      <EditButton />
    </Datagrid>
  </List>
);

export const UserCreate = props => (
  <Create {...props} title={<UserTitle label="Create new user" />}>
    <SimpleForm>
      <TextInput source="username" />
      <TextInput source="firstName" label="First Name" />
      <TextInput source="lastName" label="Last Name" />
      <PasswordInput source="password" />
      <BooleanInput source="admin" />
    </SimpleForm>
  </Create>
);

const repeatedPasswordValidator = (value, allValues) => {
  if (value !== allValues['password']) {
    return 'Password is not matching';
  }
  return undefined;
};
export const UserEdit = props => (
  <Edit {...props} mutationMode="pessimistic" title={<UserTitle label="Editing" />}>
    <SimpleForm>
      <TextField source="username" />
      <TextInput source="firstName" label="First Name" />
      <TextInput source="lastName" label="Last Name" />
      <PasswordInput source="password" />
      <PasswordInput source="password2" label="Repeat password" validate={repeatedPasswordValidator} />
      <BooleanInput source="admin" />
    </SimpleForm>
  </Edit>
);


export const UserShow = props => (
  <Show {...props} title={<UserTitle />}>
    <SimpleShowLayout>
      <TextField source="username" />
      <TextField source="firstName" label="First Name" />
      <TextField source="lastName" label="Last Name" />
      <DateField showTime source="createdAt" label="Created At" />
      <DateField showTime source="updatedAt" label="Updated At" />
      <BooleanField source="admin" />
    </SimpleShowLayout>
  </Show>
);
export const UsersComponents = {
  list: UserList,
  create: UserCreate,
  show: UserShow,
  edit: UserEdit,
  icon: PeopleIcon,
};
