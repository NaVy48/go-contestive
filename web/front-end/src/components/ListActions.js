import EventIcon from '@material-ui/icons/Event';
import React, { cloneElement } from 'react';
import { Button, CreateButton, sanitizeListRestProps, TopToolbar } from 'react-admin';

export const ListActions = ({
  currentSort,
  className,
  resource,
  filters,
  displayedFilters,
  filterValues,
  permanentFilter,
  hasCreate, // you can hide CreateButton if hasCreate = false
  basePath,
  selectedIds,
  onUnselectItems,
  showFilter,
  maxResults,
  total,
  ...rest
}) => (
  <TopToolbar className={className} {...sanitizeListRestProps(rest)}>
    {filters &&
      cloneElement(filters, {
        resource,
        showFilter,
        displayedFilters,
        filterValues,
        context: 'button',
      })}
    {!!hasCreate && <CreateButton basePath={basePath} />}
    <Button
      onClick={() => {
        alert('Your custom action');
      }}
      label="Show calendar"
    >
      <EventIcon />
    </Button>
  </TopToolbar>
);

ListActions.defaultProps = {
  selectedIds: [],
  onUnselectItems: () => null,
};
