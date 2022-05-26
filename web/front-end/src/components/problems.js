import AssignmentIcon from '@material-ui/icons/Assignment';
import React from 'react';
import {
  Button,
  Create,
  Datagrid,
  DateField,
  Edit,
  FileField,
  FileInput,
  List,
  NumberField,
  ReferenceField,
  ReferenceManyField,
  Show,
  ShowButton,
  SimpleForm,
  Tab,
  TabbedShowLayout,
  TextField,
  useRecordContext,
} from 'react-admin';
import styled from 'styled-components/macro'

const styles = {
  button: {
    marginTop: '1em',
  },
};

export const ProblemList = (props) => {
  console.log("problem list props", props)
  return (
    <List {...props} exporter={false} bulkActionButtons={false}>
      <Datagrid rowClick="show">
        <TextField source="id" />
        <ReferenceField source="authorId" reference="users">
          <TextField source="username" />
        </ReferenceField>
        <TextField source="name" />
        <TextField source="revision" label="Latest revision" />
        <TextField source="title" />
        <NumberField source="memoryLimit" />
        <NumberField source="timeLimit" />
      </Datagrid>
    </List>
  );
}

const resizeIframe = (event) => {
  event.target.style.height = event.target.contentWindow.document.documentElement.scrollHeight + 'px';
}

const HtmlField = (props) => {
  const { source } = props;
  const record = useRecordContext(props);

  return (
    <iframe srcDoc={record[source]} width="100%" frameBorder="0" onLoad={resizeIframe}></iframe>
  );
}

export const ProblemShow = props => {
  const isAdmin = props.permission === "admin"
  return (
    <Show {...props}>
      <TabbedShowLayout>
        <Tab label="summary">
          {isAdmin && (
            <>
              <TextField source="id" />
              <ReferenceField source="authorId" reference="users">
                <TextField source="username" />
              </ReferenceField>
            </>
          )}
          <TextField source="title" />
          {isAdmin && (<TextField source="revision" label="Latest revision" />)}
          <NumberField source="memoryLimit" />
          <NumberField source="timeLimit" />
          {isAdmin && (<DateField source="createdAt" />)}
          {isAdmin && (<DateField source="updatedAt" />)}
        </Tab>

        <Tab label="Statement" path="statement">
          <HtmlField source="statmentHtml" />
        </Tab>
      </TabbedShowLayout>
    </Show>
  )
};


export const ProblemUserShow = props => (
  <Show {...props}>
    <TabbedShowLayout>
      <Tab label="summary">
        <TextField source="title" />
        <TextField source="revision" label="Latest revision" />
        <NumberField source="memoryLimit" />
        <NumberField source="timeLimit" />
        <DateField source="createdAt" />
        <DateField source="updatedAt" />
      </Tab>

      <Tab label="Statement" path="statement">
        <HtmlField source="statmentHtml" />
      </Tab>
    </TabbedShowLayout>
  </Show>
);


const stopPropagation = (e) => e.stopPropagation()
const ExternalLink = ({ children, href }) => {
  return <a href={href} target="_blank" rel="noreferrer" onClick={stopPropagation}>{children}</a>
}

const StyledFileInput = styled(FileInput).attrs(props => ({
  classes: { dropZone: "drop-zone" }
}))`
.drop-zone {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 150px;
}
`;

export const ProblemsCreate = props => (
  <Create {...props}>
    <SimpleForm>
      <StyledFileInput
        source="package"
        accept=".zip"
        maxSize="20000000"
        placeholder={
          <span>Upload problem package from <ExternalLink href="https://polygon.codeforces.com/problems">polygon.codeforces.com</ExternalLink></span>
        }>
        <FileField source="package" title="title" />
      </StyledFileInput>
    </SimpleForm>
  </Create >
);


export const ProblemsEdit = props => (
  <Edit {...props} mutationMode="pessimistic">
    <SimpleForm>
      <TextField source="name" />
      <TextField source="externalURL" />
      <TextField source="revision" />
      <StyledFileInput
        source="package"
        accept=".zip"
        maxSize="20000000"
        placeholder={
          <span>Upload problem package from <ExternalLink href="https://polygon.codeforces.com/problems">polygon.codeforces.com</ExternalLink></span>
        }>
        <FileField source="package" title="title" />
      </StyledFileInput>
    </SimpleForm>
  </Edit>
);

export const ProblemsComponents = {
  list: ProblemList,
  create: ProblemsCreate,
  show: ProblemShow,
  edit: ProblemsEdit,
  icon: AssignmentIcon,
  export: false,
};

export const UserProblemsComponents = {
  show: ProblemShow,
  icon: AssignmentIcon,
};
