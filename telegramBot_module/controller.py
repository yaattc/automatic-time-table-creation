from sqlalchemy.orm import Session

from modules.schedule.classes import Lesson, Group
from modules.core import source as core


@core.db_write
def delete_all_lessons(session):
    """
    Delete all lessons from table to parse new ones

    :param session: sqlalchemy session from decorator
    """
    session.query(Lesson).delete()


@core.db_write
def insert_lesson(session, group, subject, teacher, day, start, end, room):
    """
    Insert new lesson with given parameters

    :param session: sqlalchemy session from decorator
    :param group: string
    :param subject: string
    :param teacher: string
    :param day: integer
    :param start: string
    :param end: string
    :param room: string
    """
    session.add(Lesson(group, subject, teacher, day, start, end, room))
