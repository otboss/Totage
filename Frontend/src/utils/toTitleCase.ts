export const toTitleCase = (str: string) => { 
  return str.replace(/_/g, ' ').toLowerCase().replace(/\b\w/g, char => char.toUpperCase()); 
};